ace.define("ace/mode/sql_highlight_rules",["require","exports","module","ace/lib/oop","ace/mode/text_highlight_rules"], function(require, exports, module) {
"use strict";

var oop = require("../lib/oop");
var TextHighlightRules = require("./text_highlight_rules").TextHighlightRules;

var SqlHighlightRules = function() {

    var keywords = (
        // all
        "cb_info|chrome_extensions|cpuid|etc_hosts|etc_protocols|etc_services|example|hash|" +
        "interface_addresses|interface_details|kernel_info|listening_ports|os_version|platform_info|" +
        "processes|process_open_sockets|system_info|uptime|users|" +
        // darwin
        "ad_config|alf_exceptions|alf_explicit_auths|alf_services|alf|app_schemes|apps|asl|" +
        "authorization_mechanisms|authorizations|browser_plugins|certificates|crashes|disk_events|" +
        "extended_attributes|fan_speed_sensors|homebrew_packages|iokit_devicetree|iokit_registry|" +
        "kernel_extensions|kernel_panics|keychain_acls|keychain_items|launchd_overrides|launchd|" +
        "managed_policies|nfs_shares|nvram|package_bom|package_install_history|package_receipts|" +
        "power_sensors|preferences|process_file_events|quicklock_cache|safari_extensions|sandboxes|" +
        "signature|sip_config|smc_keys|startup_items|temperature_sensors|wifi_networks|wifi_scan|" +
        "wifi_status|xprotect_entries|xprotect_meta|xprotect_reports|" +
        // linux
        "apt_sources|cpu_time|deb_packages|iptables|kernel_integrity|kernel_modules|memory_info|" +
        "memory_map|msr|portage_keywords|portage_packages|portage_use|rpm_package_files|rpm_packages|" +
        "shared_memory|socket_events|syslog|user_events|" +
        // posix
        "cpi_tables|arp_cache|augeas|authorized_keys|block_devices|crontab|device_file|device_hash|" +
        "device_partitions|disk_encryption|dns_resolvers|file_events|firefox_addons|groups|" +
        "hardware_events|known_hosts|last|logged_in_users|magic|mounts|opera_extensions|pci_devices|" +
        "process_envs|process_events|process_memory_map|process_open_files|routes|shell_history|" +
        "smbios_tables|ssh_keys|sudoers|suid_bin|system_controls|usb_devices|user_groups|yara_events|yara|" +
        // utility
        "file|osquery_events|osquery_extensions|osquery_flags|osquery_info|osquery_packs|osquery_registry|" +
        "osquery_schedule|time|" +
        // windows
        "appcompat_shims|drivers|patches|programs|registry|services|shared_resources|wmi_cli_event_consumers|" +
        "wmi_event_filters|wmi_filter_consumer_binding|wmi_script_event_consumers|" +
        // sql
        "select|insert|update|delete|from|where|and|or|group|by|order|limit|offset|having|as|case|" +
        "when|else|end|type|left|right|join|on|outer|desc|asc|union"
    );

    var builtinConstants = (
        "true|false|null"
    );

    var builtinFunctions = (
        "count|min|max|avg|sum|rank|now|coalesce"
    );

    var keywordMapper = this.createKeywordMapper({
        "support.function": builtinFunctions,
        "keyword": keywords,
        "constant.language": builtinConstants
    }, "identifier", true);

    this.$rules = {
        "start" : [ {
            token : "comment",
            regex : "--.*$"
        },  {
            token : "comment",
            start : "/\\*",
            end : "\\*/"
        }, {
            token : "string",           // " string
            regex : '".*?"'
        }, {
            token : "string",           // ' string
            regex : "'.*?'"
        }, {
            token : "constant.numeric", // float
            regex : "[+-]?\\d+(?:(?:\\.\\d*)?(?:[eE][+-]?\\d+)?)?\\b"
        }, {
            token : keywordMapper,
            regex : "[a-zA-Z_$][a-zA-Z0-9_$]*\\b"
        }, {
            token : "keyword.operator",
            regex : "\\+|\\-|\\/|\\/\\/|%|<@>|@>|<@|&|\\^|~|<|>|<=|=>|==|!=|<>|="
        }, {
            token : "paren.lparen",
            regex : "[\\(]"
        }, {
            token : "paren.rparen",
            regex : "[\\)]"
        }, {
            token : "text",
            regex : "\\s+"
        } ]
    };
    this.normalizeRules();
};

oop.inherits(SqlHighlightRules, TextHighlightRules);

exports.SqlHighlightRules = SqlHighlightRules;
});

ace.define("ace/mode/sql",["require","exports","module","ace/lib/oop","ace/mode/text","ace/mode/sql_highlight_rules","ace/range"], function(require, exports, module) {
"use strict";

var oop = require("../lib/oop");
var TextMode = require("./text").Mode;
var SqlHighlightRules = require("./sql_highlight_rules").SqlHighlightRules;
var Range = require("../range").Range;

var Mode = function() {
    this.HighlightRules = SqlHighlightRules;
};
oop.inherits(Mode, TextMode);

(function() {

    this.lineCommentStart = "--";

    this.$id = "ace/mode/sql";
}).call(Mode.prototype);

exports.Mode = Mode;

});
