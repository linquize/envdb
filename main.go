package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mephux/envdb/envdb"

	"github.com/howeyc/gopass"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app   = kingpin.New(envdb.Name, "The Environment Database - Ask your environment questions")
	debug = app.Flag("debug", "Enable debug logging.").Short('v').Bool()
	dev   = app.Flag("dev", "Enable dev mode.").Bool()
	quiet = app.Flag("quiet", "Remove all output logging.").Short('q').Bool()

	server        = app.Command("server", "Start the tcp server for node connections.")
	serverCommand = server.Arg("command", "Daemon command. (start,status,stop)").String()
	// serverConfig = server.Flag("config", "Server configuration file.").File()
	serverPort = server.Flag("port", "Port for the server to listen on.").
			Short('p').PlaceHolder(fmt.Sprintf("%d", envdb.DefaultServerPort)).Int()

	serverWebPort = server.Flag("http-port", "Port for the web server to listen on.").
			Short('P').PlaceHolder(fmt.Sprintf("%d", envdb.DefaultWebServerPort)).Int()

	node = app.Command("node", "Register a new node.")
	// clientConfig = client.Flag("config", "Client configuration file.").File()
	nodeName = node.Arg("node-name", "A name used to uniquely identify this node.").Required().String()

	nodeServer = node.Flag("server", "Address for server to connect to.").
			Short('s').PlaceHolder("127.0.0.1").Required().String()

	nodePort = node.Flag("port", "Port to use for connection.").Short('p').Int()

	users      = app.Command("users", "User Management (Default lists all users).")
	addUser    = users.Flag("add", "Add a new user.").Short('a').Bool()
	removeUser = users.Flag("remove", "Remove user by email.").Short('r').PlaceHolder("email").String()

	// Log Global logger
	Log *envdb.Logger

	// Build holds the git commit that was compiled.
	// This will be filled in by the compiler.
	Build string
)

func initLogger() {
	Log = envdb.NewLogger()

	Log.Prefix = envdb.Name

	if *debug {
		Log.SetLevel(envdb.DebugLevel)
	} else {
		Log.SetLevel(envdb.InfoLevel)
	}

	if *dev {
		envdb.DevMode = true
		Log.SetLevel(envdb.DebugLevel)
		Log.Info("DEBUG MODE ENABLED.")
	} else {
		envdb.DevMode = false
	}

	if *quiet {
		Log.SetLevel(envdb.FatalLevel)
	}

	envdb.Log = Log
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Version(fmt.Sprintf("%s-%s", envdb.Version, Build))

	args, err := app.Parse(os.Args[1:])

	initLogger()

	switch kingpin.MustParse(args, err) {

	case users.FullCommand():
		serverSetup(false)

		if *addUser {
			addDBUser()
			return
		}

		if len(*removeUser) > 0 {
			if user, err := envdb.FindUserByEmail(*removeUser); err != nil {
				Log.Fatal(err)
			} else {
				if err := user.Delete(); err != nil {
					Log.Fatal(err)
				}
			}

			Log.Info("User removed successfully.")
			return
		}

		users, err := envdb.FindAllUsers()

		if err != nil {
			Log.Fatal(err)
		}

		fmt.Println("Listing Users: ")
		for _, user := range users {
			fmt.Printf("  * %s (%s)\n", user.Name, user.Email)
		}

	case server.FullCommand():

		serverSetup(true)

	case node.FullCommand():

		output, err := exec.Command("whoami").Output()

		if err != nil {
			Log.Fatalf("Error: %s", err)
			os.Exit(-1)
		}

		if strings.Trim(string(output), "\n") != "root" {
			Log.Fatal("You must run the node client as root.")
			os.Exit(-1)
		}

		var clntPort = envdb.DefaultServerPort

		if *nodePort != 0 {
			clntPort = *nodePort
		}

		var c = envdb.Node{
			Name:       *nodeName,
			Host:       *nodeServer,
			Port:       clntPort,
			RetryCount: 50,
		}

		config, err := envdb.NewNodeConfig()

		if err != nil {
			Log.Fatal(err)
		}

		c.Config = config

		if err := c.Run(); err != nil {
			Log.Error(err)
		}

	default:
		app.Usage([]string{})
	}

}

func serverSetup(start bool) {
	var svrPort = envdb.DefaultServerPort
	var svrWebPort = envdb.DefaultWebServerPort

	if *serverPort != 0 {
		svrPort = *serverPort
	}

	if *serverWebPort != 0 {
		svrWebPort = *serverWebPort
	}

	svr, err := envdb.NewServer(svrPort)

	if err != nil {
		Log.Fatal(err)
	}

	if !start {
		return
	}

	if len(*serverCommand) <= 0 {
		if err := svr.Run(svrWebPort); err != nil {
			Log.Error(err)
		}
	} else {

		switch *serverCommand {
		case "start":
			svr.Config.Daemon.StartServer(svr, svrWebPort)
			break
		case "stop":
			svr.Config.Daemon.Stop()
			break
		case "status":
			svr.Config.Daemon.Status()
			break
		default:
			{
				Log.Fatalf("Error: Unknown Command %s.", *serverCommand)
			}
		}

	}

}

func ask(reader *bufio.Reader, question string) string {
	fmt.Print(question)

	value, _ := reader.ReadString('\n')
	trim := strings.Trim(value, "\n")

	if len(trim) <= 0 {
		Log.Fatalf("value cannot be blank.")
	}

	return trim
}

func addDBUser() {
	reader := bufio.NewReader(os.Stdin)

	name := ask(reader, "Name: ")

	email := ask(reader, "Email ")

	if !envdb.IsEmail(email) {
		Log.Fatalf("%s is not a valid email address.", email)
	}

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()

	if err != nil {
		Log.Fatal("password stdin error")
	}

	fmt.Print("Confirm: ")
	cpass, err := gopass.GetPasswd()

	if err != nil {
		Log.Fatal("password stdin error")
	}

	if string(pass) != string(cpass) {
		Log.Fatal("Password and confirm do not match.")
	}

	user := &envdb.UserDb{
		Name:     name,
		Email:    email,
		Password: string(pass),
	}

	err = envdb.CreateUser(user)

	if err != nil {
		Log.Fatal(err)
	}

	Log.Info("User created successfully.")
}
