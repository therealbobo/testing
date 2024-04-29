// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package rules

import (
	"github.com/falcosecurity/testing/pkg/run"
)

// LegacyFalcoRules_v1_0_1 is a local copy of:
// https://github.com/falcosecurity/rules/blob/falco-rules-1.0.1/rules/falco_rules.yaml
//
// It is intended to be used with legacy tests.
// Please note that legacy tests are designed to test Falco features, not the ruleset.
var LegacyFalcoRules_v1_0_1 = run.NewStringFileAccessor(
	"falco_rules.yaml",
	`
    #
    # Copyright (C) 2023 The Falco Authors.
    #
    #
    # Licensed under the Apache License, Version 2.0 (the "License");
    # you may not use this file except in compliance with the License.
    # You may obtain a copy of the License at
    #
    #     http://www.apache.org/licenses/LICENSE-2.0
    #
    # Unless required by applicable law or agreed to in writing, software
    # distributed under the License is distributed on an "AS IS" BASIS,
    # WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    # See the License for the specific language governing permissions and
    # limitations under the License.
    #
    
    # Starting with version 8, the Falco engine supports exceptions.
    # However the Falco rules file does not use them by default.
    - required_engine_version: 17
    
    # Currently disabled as read/write are ignored syscalls. The nearly
    # similar open_write/open_read check for files being opened for
    # reading/writing.
    # - macro: write
    #   condition: (syscall.type=write and fd.type in (file, directory))
    # - macro: read
    #   condition: (syscall.type=read and evt.dir=> and fd.type in (file, directory))
    
    # Information about rules tags and fields can be found here: https://falco.org/docs/rules/#tags-for-current-falco-ruleset
    # `+"`"+`tags`+"`"+` fields also include information about the type of workload inspection, Mitre Attack killchain phases and Mitre TTP code(s)
    # Mitre Attack References:
    # [1] https://attack.mitre.org/tactics/enterprise/
    # [2] https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/enterprise-attack.json
    
    - macro: open_write
      condition: (evt.type in (open,openat,openat2) and evt.is_open_write=true and fd.typechar='f' and fd.num>=0)
    
    - macro: open_read
      condition: (evt.type in (open,openat,openat2) and evt.is_open_read=true and fd.typechar='f' and fd.num>=0)
    
    - macro: open_directory
      condition: (evt.type in (open,openat,openat2) and evt.is_open_read=true and fd.typechar='d' and fd.num>=0)
    
    # Failed file open attempts, useful to detect threat actors making mistakes
    # https://man7.org/linux/man-pages/man3/errno.3.html
    # evt.res=ENOENT - No such file or directory
    # evt.res=EACCESS - Permission denied
    - macro: open_file_failed
      condition: (evt.type in (open,openat,openat2) and fd.typechar='f' and fd.num=-1 and evt.res startswith E)
    
    - macro: never_true
      condition: (evt.num=0)
    
    - macro: always_true
      condition: (evt.num>=0)
    
    # In some cases, such as dropped system call events, information about
    # the process name may be missing. For some rules that really depend
    # on the identity of the process performing an action such as opening
    # a file, etc., we require that the process name be known.
    - macro: proc_name_exists
      condition: (not proc.name in ("<NA>","N/A"))
    
    - macro: rename
      condition: (evt.type in (rename, renameat, renameat2))
    
    - macro: mkdir
      condition: (evt.type in (mkdir, mkdirat))
    
    - macro: remove
      condition: (evt.type in (rmdir, unlink, unlinkat))
    
    - macro: modify
      condition: (rename or remove)
    
    # %evt.arg.flags available for evt.dir=>, but only for umount2
    # %evt.arg.name is path and available for evt.dir=<
    # - macro: umount
    #   condition: (evt.type in (umount, umount2))
    
    - macro: spawned_process
      condition: (evt.type in (execve, execveat) and evt.dir=<)
    
    - macro: create_symlink
      condition: (evt.type in (symlink, symlinkat) and evt.dir=<)
    
    - macro: create_hardlink
      condition: (evt.type in (link, linkat) and evt.dir=<)
    
    - macro: chmod
      condition: (evt.type in (chmod, fchmod, fchmodat) and evt.dir=<)
    
    - macro: kernel_module_load
      condition: (evt.type in (init_module, finit_module) and evt.dir=<)
    
    # File categories
    - macro: bin_dir
      condition: (fd.directory in (/bin, /sbin, /usr/bin, /usr/sbin))
    
    - macro: bin_dir_mkdir
      condition: >
         (evt.arg.path startswith /bin/ or
         evt.arg.path startswith /sbin/ or
         evt.arg.path startswith /usr/bin/ or
         evt.arg.path startswith /usr/sbin/)
    
    - macro: bin_dir_rename
      condition: >
         (evt.arg.path startswith /bin/ or
         evt.arg.path startswith /sbin/ or
         evt.arg.path startswith /usr/bin/ or
         evt.arg.path startswith /usr/sbin/ or
         evt.arg.name startswith /bin/ or
         evt.arg.name startswith /sbin/ or
         evt.arg.name startswith /usr/bin/ or
         evt.arg.name startswith /usr/sbin/ or
         evt.arg.oldpath startswith /bin/ or
         evt.arg.oldpath startswith /sbin/ or
         evt.arg.oldpath startswith /usr/bin/ or
         evt.arg.oldpath startswith /usr/sbin/ or
         evt.arg.newpath startswith /bin/ or
         evt.arg.newpath startswith /sbin/ or
         evt.arg.newpath startswith /usr/bin/ or
         evt.arg.newpath startswith /usr/sbin/)
    
    - macro: etc_dir
      condition: (fd.name startswith /etc/)
    
    # This detects writes immediately below / or any write anywhere below /root
    - macro: root_dir
      condition: (fd.directory=/ or fd.name startswith /root/)
    
    - list: shell_binaries
      items: [ash, bash, csh, ksh, sh, tcsh, zsh, dash]
    
    - list: ssh_binaries
      items: [
        sshd, sftp-server, ssh-agent,
        ssh, scp, sftp,
        ssh-keygen, ssh-keysign, ssh-keyscan, ssh-add
        ]
    
    - list: shell_mgmt_binaries
      items: [add-shell, remove-shell]
    
    - macro: shell_procs
      condition: proc.name in (shell_binaries)
    
    - list: coreutils_binaries
      items: [
        truncate, sha1sum, numfmt, fmt, fold, uniq, cut, who,
        groups, csplit, sort, expand, printf, printenv, unlink, tee, chcon, stat,
        basename, split, nice, "yes", whoami, sha224sum, hostid, users, stdbuf,
        base64, unexpand, cksum, od, paste, nproc, pathchk, sha256sum, wc, test,
        comm, arch, du, factor, sha512sum, md5sum, tr, runcon, env, dirname,
        tsort, join, shuf, install, logname, pinky, nohup, expr, pr, tty, timeout,
        tail, "[", seq, sha384sum, nl, head, id, mkfifo, sum, dircolors, ptx, shred,
        tac, link, chroot, vdir, chown, touch, ls, dd, uname, "true", pwd, date,
        chgrp, chmod, mktemp, cat, mknod, sync, ln, "false", rm, mv, cp, echo,
        readlink, sleep, stty, mkdir, df, dir, rmdir, touch
        ]
    
    # dpkg -L login | grep bin | xargs ls -ld | grep -v '^d' | awk '{print $9}' | xargs -L 1 basename | tr "\\n" ","
    - list: login_binaries
      items: [
        login, systemd, '"(systemd)"', systemd-logind, su,
        nologin, faillog, lastlog, newgrp, sg
        ]
    
    # dpkg -L passwd | grep bin | xargs ls -ld | grep -v '^d' | awk '{print $9}' | xargs -L 1 basename | tr "\\n" ","
    - list: passwd_binaries
      items: [
        shadowconfig, grpck, pwunconv, grpconv, pwck,
        groupmod, vipw, pwconv, useradd, newusers, cppw, chpasswd, usermod,
        groupadd, groupdel, grpunconv, chgpasswd, userdel, chage, chsh,
        gpasswd, chfn, expiry, passwd, vigr, cpgr, adduser, addgroup, deluser, delgroup
        ]
    
    # repoquery -l shadow-utils | grep bin | xargs ls -ld | grep -v '^d' |
    #     awk '{print $9}' | xargs -L 1 basename | tr "\\n" ","
    - list: shadowutils_binaries
      items: [
        chage, gpasswd, lastlog, newgrp, sg, adduser, deluser, chpasswd,
        groupadd, groupdel, addgroup, delgroup, groupmems, groupmod, grpck, grpconv, grpunconv,
        newusers, pwck, pwconv, pwunconv, useradd, userdel, usermod, vigr, vipw, unix_chkpwd
        ]
    
    - list: sysdigcloud_binaries
      items: [setup-backend, dragent, sdchecks]
    
    - list: k8s_binaries
      items: [hyperkube, skydns, kube2sky, exechealthz, weave-net, loopback, bridge, openshift-sdn, openshift]
    
    - list: lxd_binaries
      items: [lxd, lxcfs]
    
    - list: http_server_binaries
      items: [nginx, httpd, httpd-foregroun, lighttpd, apache, apache2]
    
    - list: db_server_binaries
      items: [mysqld, postgres, sqlplus]
    
    - list: postgres_mgmt_binaries
      items: [pg_dumpall, pg_ctl, pg_lsclusters, pg_ctlcluster]
    
    - list: nosql_server_binaries
      items: [couchdb, memcached, redis-server, rabbitmq-server, mongod]
    
    - list: gitlab_binaries
      items: [gitlab-shell, gitlab-mon, gitlab-runner-b, git]
    
    - list: interpreted_binaries
      items: [lua, node, perl, perl5, perl6, php, python, python2, python3, ruby, tcl]
    
    - macro: interpreted_procs
      condition: >
        (proc.name in (interpreted_binaries))
    
    - macro: server_procs
      condition: proc.name in (http_server_binaries, db_server_binaries, docker_binaries, sshd)
    
    # The explicit quotes are needed to avoid the - characters being
    # interpreted by the filter expression.
    - list: rpm_binaries
      items: [dnf, dnf-automatic, rpm, rpmkey, yum, '"75-system-updat"', rhsmcertd-worke, rhsmcertd, subscription-ma,
              repoquery, rpmkeys, rpmq, yum-cron, yum-config-mana, yum-debug-dump,
              abrt-action-sav, rpmdb_stat, microdnf, rhn_check, yumdb]
    
    - list: openscap_rpm_binaries
      items: [probe_rpminfo, probe_rpmverify, probe_rpmverifyfile, probe_rpmverifypackage]
    
    - macro: rpm_procs
      condition: (proc.name in (rpm_binaries, openscap_rpm_binaries) or proc.name in (salt-call, salt-minion))
    
    - list: deb_binaries
      items: [dpkg, dpkg-preconfigu, dpkg-reconfigur, dpkg-divert, apt, apt-get, aptitude,
        frontend, preinst, add-apt-reposit, apt-auto-remova, apt-key,
        apt-listchanges, unattended-upgr, apt-add-reposit, apt-cache, apt.systemd.dai
        ]
    - list: python_package_managers
      items: [pip, pip3, conda]
    
    # The truncated dpkg-preconfigu is intentional, process names are
    # truncated at the falcosecurity-libs level.
    - list: package_mgmt_binaries
      items: [rpm_binaries, deb_binaries, update-alternat, gem, npm, python_package_managers, sane-utils.post, alternatives, chef-client, apk, snapd]
    
    - macro: package_mgmt_procs
      condition: proc.name in (package_mgmt_binaries)
    
    - macro: package_mgmt_ancestor_procs
      condition: proc.pname in (package_mgmt_binaries) or
                 proc.aname[2] in (package_mgmt_binaries) or
                 proc.aname[3] in (package_mgmt_binaries) or
                 proc.aname[4] in (package_mgmt_binaries)
    
    - macro: coreos_write_ssh_dir
      condition: (proc.name=update-ssh-keys and fd.name startswith /home/core/.ssh)
    
    - macro: run_by_package_mgmt_binaries
      condition: proc.aname in (package_mgmt_binaries, needrestart)
    
    - list: ssl_mgmt_binaries
      items: [ca-certificates]
    
    - list: dhcp_binaries
      items: [dhclient, dhclient-script, 11-dhclient]
    
    # A canonical set of processes that run other programs with different
    # privileges or as a different user.
    - list: userexec_binaries
      items: [sudo, su, suexec, critical-stack, dzdo]
    
    - list: known_setuid_binaries
      items: [
        sshd, dbus-daemon-lau, ping, ping6, critical-stack-, pmmcli,
        filemng, PassengerAgent, bwrap, osdetect, nginxmng, sw-engine-fpm,
        start-stop-daem
        ]
    
    - list: user_mgmt_binaries
      items: [login_binaries, passwd_binaries, shadowutils_binaries]
    
    - list: dev_creation_binaries
      items: [blkid, rename_device, update_engine, sgdisk]
    
    - list: hids_binaries
      items: [aide, aide.wrapper, update-aide.con, logcheck, syslog-summary, osqueryd, ossec-syscheckd]
    
    - list: vpn_binaries
      items: [openvpn]
    
    - list: nomachine_binaries
      items: [nxexec, nxnode.bin, nxserver.bin, nxclient.bin]
    
    - macro: system_procs
      condition: proc.name in (coreutils_binaries, user_mgmt_binaries)
    
    - list: mail_binaries
      items: [
        sendmail, sendmail-msp, postfix, procmail, exim4,
        pickup, showq, mailq, dovecot, imap-login, imap,
        mailmng-core, pop3-login, dovecot-lda, pop3
        ]
    
    - list: mail_config_binaries
      items: [
        update_conf, parse_mc, makemap_hash, newaliases, update_mk, update_tlsm4,
        update_db, update_mc, ssmtp.postinst, mailq, postalias, postfix.config.,
        postfix.config, postfix-script, postconf
        ]
    
    - list: sensitive_file_names
      items: [/etc/shadow, /etc/sudoers, /etc/pam.conf, /etc/security/pwquality.conf]
    
    - list: sensitive_directory_names
      items: [/, /etc, /etc/, /root, /root/]
    
    - macro: sensitive_files
      condition: >
        fd.name startswith /etc and
        (fd.name in (sensitive_file_names)
         or fd.directory in (/etc/sudoers.d, /etc/pam.d))
    
    # Indicates that the process is new. Currently detected using time
    # since process was started, using a threshold of 5 seconds.
    - macro: proc_is_new
      condition: proc.duration <= 5000000000
    
    # Network
    - macro: inbound
      condition: >
        (((evt.type in (accept,accept4,listen) and evt.dir=<) or
          (evt.type in (recvfrom,recvmsg) and evt.dir=< and
           fd.l4proto != tcp and fd.connected=false and fd.name_changed=true)) and
         (fd.typechar = 4 or fd.typechar = 6) and
         (fd.ip != "0.0.0.0" and fd.net != "127.0.0.0/8") and
         (evt.rawres >= 0 or evt.res = EINPROGRESS))
    
    # RFC1918 addresses were assigned for private network usage
    - list: rfc_1918_addresses
      items: ['"10.0.0.0/8"', '"172.16.0.0/12"', '"192.168.0.0/16"']
    
    - macro: outbound
      condition: >
        (((evt.type = connect and evt.dir=<) or
          (evt.type in (sendto,sendmsg) and evt.dir=< and
           fd.l4proto != tcp and fd.connected=false and fd.name_changed=true)) and
         (fd.typechar = 4 or fd.typechar = 6) and
         (fd.ip != "0.0.0.0" and fd.net != "127.0.0.0/8" and not fd.snet in (rfc_1918_addresses)) and
         (evt.rawres >= 0 or evt.res = EINPROGRESS))
    
    # Very similar to inbound/outbound, but combines the tests together
    # for efficiency.
    - macro: inbound_outbound
      condition: >
        ((((evt.type in (accept,accept4,listen,connect) and evt.dir=<)) and
         (fd.typechar = 4 or fd.typechar = 6)) and
         (fd.ip != "0.0.0.0" and fd.net != "127.0.0.0/8") and
         (evt.rawres >= 0 or evt.res = EINPROGRESS))
    
    - macro: ssh_port
      condition: fd.sport=22
    
    # In a local/user rules file, you could override this macro to
    # enumerate the servers for which ssh connections are allowed. For
    # example, you might have a ssh gateway host for which ssh connections
    # are allowed.
    #
    # In the main falco rules file, there isn't any way to know the
    # specific hosts for which ssh access is allowed, so this macro just
    # repeats ssh_port, which effectively allows ssh from all hosts. In
    # the overridden macro, the condition would look something like
    # "fd.sip="a.b.c.d" or fd.sip="e.f.g.h" or ..."
    - macro: allowed_ssh_hosts
      condition: ssh_port
    
    - rule: Disallowed SSH Connection
      desc: Detect any new ssh connection to a host other than those in an allowed group of hosts
      condition: (inbound_outbound) and ssh_port and not allowed_ssh_hosts
      enabled: false
      output: Disallowed SSH Connection (command=%proc.cmdline pid=%proc.pid connection=%fd.name user=%user.name user_loginuid=%user.loginuid container_id=%container.id image=%container.image.repository)
      priority: NOTICE
      tags: [host, container, network, mitre_command_and_control, mitre_lateral_movement, T1021.004]
    
    # These rules and supporting macros are more of an example for how to
    # use the fd.*ip and fd.*ip.name fields to match connection
    # information against ips, netmasks, and complete domain names.
    #
    # To use this rule, you should enable it and
    # populate allowed_{source,destination}_{ipaddrs,networks,domains} with the
    # values that make sense for your environment.
    
    # Note that this can be either individual IPs or netmasks
    - list: allowed_outbound_destination_ipaddrs
      items: ['"127.0.0.1"', '"8.8.8.8"']
    
    - list: allowed_outbound_destination_networks
      items: ['"127.0.0.1/8"']
    
    - list: allowed_outbound_destination_domains
      items: [google.com, www.yahoo.com]
    
    - rule: Unexpected outbound connection destination
      desc: Detect any outbound connection to a destination outside of an allowed set of ips, networks, or domain names
      condition: >
        outbound and not
        ((fd.sip in (allowed_outbound_destination_ipaddrs)) or
         (fd.snet in (allowed_outbound_destination_networks)) or
         (fd.sip.name in (allowed_outbound_destination_domains)))
      enabled: false
      output: Disallowed outbound connection destination (command=%proc.cmdline pid=%proc.pid connection=%fd.name user=%user.name user_loginuid=%user.loginuid container_id=%container.id image=%container.image.repository)
      priority: NOTICE
      tags: [host, container, network, mitre_command_and_control, TA0011]
    
    - list: allowed_inbound_source_ipaddrs
      items: ['"127.0.0.1"']
    
    - list: allowed_inbound_source_networks
      items: ['"127.0.0.1/8"', '"10.0.0.0/8"']
    
    - list: allowed_inbound_source_domains
      items: [google.com]
    
    - rule: Unexpected inbound connection source
      desc: Detect any inbound connection from a source outside of an allowed set of ips, networks, or domain names
      condition: >
        inbound and not
        ((fd.cip in (allowed_inbound_source_ipaddrs)) or
         (fd.cnet in (allowed_inbound_source_networks)) or
         (fd.cip.name in (allowed_inbound_source_domains)))
      enabled: false
      output: Disallowed inbound connection source (command=%proc.cmdline pid=%proc.pid connection=%fd.name user=%user.name user_loginuid=%user.loginuid container_id=%container.id image=%container.image.repository)
      priority: NOTICE
      tags: [host, container, network, mitre_command_and_control, TA0011]
    
    - list: bash_config_filenames
      items: [.bashrc, .bash_profile, .bash_history, .bash_login, .bash_logout, .inputrc, .profile]
    
    - list: bash_config_files
      items: [/etc/profile, /etc/bashrc]
    
    # Covers both csh and tcsh
    - list: csh_config_filenames
      items: [.cshrc, .login, .logout, .history, .tcshrc, .cshdirs]
    
    - list: csh_config_files
      items: [/etc/csh.cshrc, /etc/csh.login]
    
    - list: zsh_config_filenames
      items: [.zshenv, .zprofile, .zshrc, .zlogin, .zlogout]
    
    - list: shell_config_filenames
      items: [bash_config_filenames, csh_config_filenames, zsh_config_filenames]
    
    - list: shell_config_files
      items: [bash_config_files, csh_config_files]
    
    - list: shell_config_directories
      items: [/etc/zsh]
    
    - macro: user_known_shell_config_modifiers
      condition: (never_true)
    
    - rule: Modify Shell Configuration File
      desc: Detect attempt to modify shell configuration files
      condition: >
        open_write and
        (fd.filename in (shell_config_filenames) or
         fd.name in (shell_config_files) or
         fd.directory in (shell_config_directories))
        and not proc.name in (shell_binaries)
        and not exe_running_docker_save
        and not user_known_shell_config_modifiers
      output: >
        a shell configuration file has been modified (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid pcmdline=%proc.pcmdline file=%fd.name container_id=%container.id image=%container.image.repository)
      priority:
        WARNING
      tags: [host, container, filesystem, mitre_persistence, T1546.004]
    
    # This rule is not enabled by default, as there are many legitimate
    # readers of shell config files.
    - rule: Read Shell Configuration File
      desc: Detect attempts to read shell configuration files by non-shell programs
      condition: >
        open_read and
        (fd.filename in (shell_config_filenames) or
         fd.name in (shell_config_files) or
         fd.directory in (shell_config_directories)) and
        (not proc.name in (shell_binaries))
      enabled: false
      output: >
        a shell configuration file was read by a non-shell program (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid file=%fd.name container_id=%container.id image=%container.image.repository)
      priority:
        WARNING
      tags: [host, container, filesystem, mitre_discovery, T1546.004]
    
    - macro: user_known_cron_jobs
      condition: (never_true)
    
    - rule: Schedule Cron Jobs
      desc: Detect cron jobs scheduled
      condition: >
        ((open_write and fd.name startswith /etc/cron) or
         (spawned_process and proc.name = "crontab")) and
        not user_known_cron_jobs
      enabled: false
      output: >
        Cron jobs were scheduled to run (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid
        file=%fd.name container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag)
      priority:
        NOTICE
      tags: [host, container, filesystem, mitre_persistence, T1053.003]
    
    # Use this to test whether the event occurred within a container.
    
    # When displaying container information in the output field, use
    # %container.info, without any leading term (file=%fd.name
    # %container.info user=%user.name user_loginuid=%user.loginuid, and not file=%fd.name
    # container=%container.info user=%user.name user_loginuid=%user.loginuid). The output will change
    # based on the context and whether or not -pk/-pm/-pc was specified on
    # the command line.
    - macro: container
      condition: (container.id != host)
    
    - macro: container_started
      condition: >
        ((evt.type = container or
         (spawned_process and proc.vpid=1)) and
         container.image.repository != incomplete)
    
    - macro: interactive
      condition: >
        ((proc.aname=sshd and proc.name != sshd) or
        proc.name=systemd-logind or proc.name=login)
    
    - list: cron_binaries
      items: [anacron, cron, crond, crontab]
    
    # https://github.com/liske/needrestart
    - list: needrestart_binaries
      items: [needrestart, 10-dpkg, 20-rpm, 30-pacman]
    
    # Possible scripts run by sshkit
    - list: sshkit_script_binaries
      items: [10_etc_sudoers., 10_passwd_group]
    
    - list: plesk_binaries
      items: [sw-engine, sw-engine-fpm, sw-engine-kv, filemng, f2bmng]
    
    # System users that should never log into a system. Consider adding your own
    # service users (e.g. 'apache' or 'mysqld') here.
    - macro: system_users
      condition: user.name in (bin, daemon, games, lp, mail, nobody, sshd, sync, uucp, www-data)
    
    - macro: httpd_writing_ssl_conf
      condition: >
        (proc.pname=run-httpd and
         (proc.cmdline startswith "sed -ri" or proc.cmdline startswith "sed -i") and
         (fd.name startswith /etc/httpd/conf.d/ or fd.name startswith /etc/httpd/conf))
    
    - macro: userhelper_writing_etc_security
      condition: (proc.name=userhelper and fd.name startswith /etc/security)
    
    - macro: ansible_running_python
      condition: (proc.name in (python, pypy, python3) and proc.cmdline contains ansible)
    
    - macro: python_running_chef
      condition: (proc.name=python and (proc.cmdline contains yum-dump.py or proc.cmdline="python /usr/bin/chef-monitor.py"))
    
    - macro: python_running_denyhosts
      condition: >
        (proc.name=python and
        (proc.cmdline contains /usr/sbin/denyhosts or
         proc.cmdline contains /usr/local/bin/denyhosts.py))
    
    # Qualys seems to run a variety of shell subprocesses, at various
    # levels. This checks at a few levels without the cost of a full
    # proc.aname, which traverses the full parent hierarchy.
    - macro: run_by_qualys
      condition: >
        (proc.pname=qualys-cloud-ag or
         proc.aname[2]=qualys-cloud-ag or
         proc.aname[3]=qualys-cloud-ag or
         proc.aname[4]=qualys-cloud-ag)
    
    - macro: run_by_sumologic_securefiles
      condition: >
        ((proc.cmdline="usermod -a -G sumologic_collector" or
          proc.cmdline="groupadd sumologic_collector") and
         (proc.pname=secureFiles.sh and proc.aname[2]=java))
    
    - macro: run_by_yum
      condition: ((proc.pname=sh and proc.aname[2]=yum) or
                  (proc.aname[2]=sh and proc.aname[3]=yum))
    
    - macro: run_by_ms_oms
      condition: >
        (proc.aname[3] startswith omsagent- or
         proc.aname[3] startswith scx-)
    
    - macro: run_by_google_accounts_daemon
      condition: >
        (proc.aname[1] startswith google_accounts or
         proc.aname[2] startswith google_accounts or
         proc.aname[3] startswith google_accounts)
    
    # Chef is similar.
    - macro: run_by_chef
      condition: (proc.aname[2]=chef_command_wr or proc.aname[3]=chef_command_wr or
                  proc.aname[2]=chef-client or proc.aname[3]=chef-client or
                  proc.name=chef-client)
    
    - macro: run_by_adclient
      condition: (proc.aname[2]=adclient or proc.aname[3]=adclient or proc.aname[4]=adclient)
    
    - macro: run_by_centrify
      condition: (proc.aname[2]=centrify or proc.aname[3]=centrify or proc.aname[4]=centrify)
    
    # Also handles running semi-indirectly via scl
    - macro: run_by_foreman
      condition: >
        (user.name=foreman and
         ((proc.pname in (rake, ruby, scl) and proc.aname[5] in (tfm-rake,tfm-ruby)) or
         (proc.pname=scl and proc.aname[2] in (tfm-rake,tfm-ruby))))
    
    - macro: java_running_sdjagent
      condition: proc.name=java and proc.cmdline contains sdjagent.jar
    
    - macro: kubelet_running_loopback
      condition: (proc.pname=kubelet and proc.name=loopback)
    
    - macro: python_mesos_marathon_scripting
      condition: (proc.pcmdline startswith "python3 /marathon-lb/marathon_lb.py")
    
    - macro: splunk_running_forwarder
      condition: (proc.pname=splunkd and proc.cmdline startswith "sh -c /opt/splunkforwarder")
    
    - macro: parent_supervise_running_multilog
      condition: (proc.name=multilog and proc.pname=supervise)
    
    - macro: supervise_writing_status
      condition: (proc.name in (supervise,svc) and fd.name startswith "/etc/sb/")
    
    - macro: pki_realm_writing_realms
      condition: (proc.cmdline startswith "bash /usr/local/lib/pki/pki-realm" and fd.name startswith /etc/pki/realms)
    
    - macro: htpasswd_writing_passwd
      condition: (proc.name=htpasswd and fd.name=/etc/nginx/.htpasswd)
    
    - macro: lvprogs_writing_conf
      condition: >
        (proc.name in (dmeventd,lvcreate,pvscan,lvs) and
         (fd.name startswith /etc/lvm/archive or
          fd.name startswith /etc/lvm/backup or
          fd.name startswith /etc/lvm/cache))
    
    - macro: ovsdb_writing_openvswitch
      condition: (proc.name=ovsdb-server and fd.directory=/etc/openvswitch)
    
    - macro: perl_running_plesk
      condition: (proc.cmdline startswith "perl /opt/psa/admin/bin/plesk_agent_manager" or
                  proc.pcmdline startswith "perl /opt/psa/admin/bin/plesk_agent_manager")
    
    - macro: perl_running_updmap
      condition: (proc.cmdline startswith "perl /usr/bin/updmap")
    
    - macro: perl_running_centrifydc
      condition: (proc.cmdline startswith "perl /usr/share/centrifydc")
    
    - macro: runuser_reading_pam
      condition: (proc.name=runuser and fd.directory=/etc/pam.d)
    
    # CIS Linux Benchmark program
    - macro: linux_bench_reading_etc_shadow
      condition: ((proc.aname[2]=linux-bench and
                   proc.name in (awk,cut,grep)) and
                  (fd.name=/etc/shadow or
                   fd.directory=/etc/pam.d))
    
    - macro: parent_ucf_writing_conf
      condition: (proc.pname=ucf and proc.aname[2]=frontend)
    
    - macro: consul_template_writing_conf
      condition: >
        ((proc.name=consul-template and fd.name startswith /etc/haproxy) or
         (proc.name=reload.sh and proc.aname[2]=consul-template and fd.name startswith /etc/ssl))
    
    - macro: countly_writing_nginx_conf
      condition: (proc.cmdline startswith "nodejs /opt/countly/bin" and fd.name startswith /etc/nginx)
    
    - list: ms_oms_binaries
      items: [omi.postinst, omsconfig.posti, scx.postinst, omsadmin.sh, omiagent]
    
    - macro: ms_oms_writing_conf
      condition: >
        ((proc.name in (omiagent,omsagent,in_heartbeat_r*,omsadmin.sh,PerformInventor,dsc_host)
           or proc.pname in (ms_oms_binaries)
           or proc.aname[2] in (ms_oms_binaries))
         and (fd.name startswith /etc/opt/omi or fd.name startswith /etc/opt/microsoft/omsagent))
    
    - macro: ms_scx_writing_conf
      condition: (proc.name in (GetLinuxOS.sh) and fd.name startswith /etc/opt/microsoft/scx)
    
    - macro: azure_scripts_writing_conf
      condition: (proc.pname startswith "bash /var/lib/waagent/" and fd.name startswith /etc/azure)
    
    - macro: azure_networkwatcher_writing_conf
      condition: (proc.name in (NetworkWatcherA) and fd.name=/etc/init.d/AzureNetworkWatcherAgent)
    
    - macro: couchdb_writing_conf
      condition: (proc.name=beam.smp and proc.cmdline contains couchdb and fd.name startswith /etc/couchdb)
    
    - macro: update_texmf_writing_conf
      condition: (proc.name=update-texmf and fd.name startswith /etc/texmf)
    
    - macro: slapadd_writing_conf
      condition: (proc.name=slapadd and fd.name startswith /etc/ldap)
    
    - macro: openldap_writing_conf
      condition: (proc.pname=run-openldap.sh and fd.name startswith /etc/openldap)
    
    - macro: ucpagent_writing_conf
      condition: (proc.name=apiserver and container.image.repository=docker/ucp-agent and fd.name=/etc/authorization_config.cfg)
    
    - macro: iscsi_writing_conf
      condition: (proc.name=iscsiadm and fd.name startswith /etc/iscsi)
    
    - macro: istio_writing_conf
      condition: (proc.name=pilot-agent and fd.name startswith /etc/istio)
    
    - macro: symantec_writing_conf
      condition: >
        ((proc.name=symcfgd and fd.name startswith /etc/symantec) or
         (proc.name=navdefutil and fd.name=/etc/symc-defutils.conf))
    
    - macro: liveupdate_writing_conf
      condition: (proc.cmdline startswith "java LiveUpdate" and fd.name in (/etc/liveupdate.conf, /etc/Product.Catalog.JavaLiveUpdate))
    
    - macro: rancher_agent
      condition: (proc.name=agent and container.image.repository contains "rancher/agent")
    
    - macro: rancher_network_manager
      condition: (proc.name=rancher-bridge and container.image.repository contains "rancher/network-manager")
    
    - macro: sosreport_writing_files
      condition: >
        (proc.name=urlgrabber-ext- and proc.aname[3]=sosreport and
         (fd.name startswith /etc/pkt/nssdb or fd.name startswith /etc/pki/nssdb))
    
    - macro: pkgmgmt_progs_writing_pki
      condition: >
        (proc.name=urlgrabber-ext- and proc.pname in (yum, yum-cron, repoquery) and
         (fd.name startswith /etc/pkt/nssdb or fd.name startswith /etc/pki/nssdb))
    
    - macro: update_ca_trust_writing_pki
      condition: (proc.pname=update-ca-trust and proc.name=trust and fd.name startswith /etc/pki)
    
    - macro: brandbot_writing_os_release
      condition: proc.name=brandbot and fd.name=/etc/os-release
    
    - macro: selinux_writing_conf
      condition: (proc.name in (semodule,genhomedircon,sefcontext_comp) and fd.name startswith /etc/selinux)
    
    - list: veritas_binaries
      items: [vxconfigd, sfcache, vxclustadm, vxdctl, vxprint, vxdmpadm, vxdisk, vxdg, vxassist, vxtune]
    
    - macro: veritas_driver_script
      condition: (proc.cmdline startswith "perl /opt/VRTSsfmh/bin/mh_driver.pl")
    
    - macro: veritas_progs
      condition: (proc.name in (veritas_binaries) or veritas_driver_script)
    
    - macro: veritas_writing_config
      condition: (veritas_progs and (fd.name startswith /etc/vx or fd.name startswith /etc/opt/VRTS or fd.name startswith /etc/vom))
    
    - macro: nginx_writing_conf
      condition: (proc.name in (nginx,nginx-ingress-c,nginx-ingress) and (fd.name startswith /etc/nginx or fd.name startswith /etc/ingress-controller))
    
    - macro: nginx_writing_certs
      condition: >
        (((proc.name=openssl and proc.pname=nginx-launch.sh) or proc.name=nginx-launch.sh) and fd.name startswith /etc/nginx/certs)
    
    - macro: chef_client_writing_conf
      condition: (proc.pcmdline startswith "chef-client /opt/gitlab" and fd.name startswith /etc/gitlab)
    
    - macro: centrify_writing_krb
      condition: (proc.name in (adjoin,addns) and fd.name startswith /etc/krb5)
    
    - macro: sssd_writing_krb
      condition: (proc.name=adcli and proc.aname[2]=sssd and fd.name startswith /etc/krb5)
    
    - macro: cockpit_writing_conf
      condition: >
        ((proc.pname=cockpit-kube-la or proc.aname[2]=cockpit-kube-la)
         and fd.name startswith /etc/cockpit)
    
    - macro: ipsec_writing_conf
      condition: (proc.name=start-ipsec.sh and fd.directory=/etc/ipsec)
    
    - macro: exe_running_docker_save
      condition: >
        proc.name = "exe"
        and (proc.cmdline contains "/var/lib/docker"
        or proc.cmdline contains "/var/run/docker")
        and proc.pname in (dockerd, docker, dockerd-current, docker-current)
    
    # Ideally we'd have a length check here as well but
    # filterchecks don't have operators like len()
    - macro: sed_temporary_file
      condition: (proc.name=sed and fd.name startswith "/etc/sed")
    
    - macro: python_running_get_pip
      condition: (proc.cmdline startswith "python get-pip.py")
    
    - macro: python_running_ms_oms
      condition: (proc.cmdline startswith "python /var/lib/waagent/")
    
    - macro: gugent_writing_guestagent_log
      condition: (proc.name=gugent and fd.name=GuestAgent.log)
    
    - macro: dse_writing_tmp
      condition: (proc.name=dse-entrypoint and fd.name=/root/tmp__)
    
    - macro: zap_writing_state
      condition: (proc.name=java and proc.cmdline contains "jar /zap" and fd.name startswith /root/.ZAP)
    
    - macro: airflow_writing_state
      condition: (proc.name=airflow and fd.name startswith /root/airflow)
    
    - macro: rpm_writing_root_rpmdb
      condition: (proc.name=rpm and fd.directory=/root/.rpmdb)
    
    - macro: maven_writing_groovy
      condition: (proc.name=java and proc.cmdline contains "classpath /usr/local/apache-maven" and fd.name startswith /root/.groovy)
    
    - macro: chef_writing_conf
      condition: (proc.name=chef-client and fd.name startswith /root/.chef)
    
    - macro: kubectl_writing_state
      condition: (proc.name in (kubectl,oc) and fd.name startswith /root/.kube)
    
    - macro: java_running_cassandra
      condition: (proc.name=java and proc.cmdline contains "cassandra.jar")
    
    - macro: cassandra_writing_state
      condition: (java_running_cassandra and fd.directory=/root/.cassandra)
    
    # Istio
    - macro: galley_writing_state
      condition: (proc.name=galley and fd.name in (known_istio_files))
    
    - list: known_istio_files
      items: [/healthready, /healthliveness]
    
    - macro: calico_writing_state
      condition: (proc.name=kube-controller and fd.name startswith /status.json and k8s.pod.name startswith calico)
    
    - macro: calico_writing_envvars
      condition: (proc.name=start_runit and fd.name startswith "/etc/envvars" and container.image.repository endswith "calico/node")
    
    - list: repository_files
      items: [sources.list]
    
    - list: repository_directories
      items: [/etc/apt/sources.list.d, /etc/yum.repos.d, /etc/apt]
    
    - macro: access_repositories
      condition: (fd.directory in (repository_directories) or
                  (fd.name pmatch (repository_directories) and
                   fd.filename in (repository_files)))
    
    - macro: modify_repositories
      condition: (evt.arg.newpath pmatch (repository_directories))
    
    - macro: user_known_update_package_registry
      condition: (never_true)
    
    - rule: Update Package Repository
      desc: Detect package repositories get updated
      condition: >
        ((open_write and access_repositories) or (modify and modify_repositories))
        and not package_mgmt_procs
        and not package_mgmt_ancestor_procs
        and not exe_running_docker_save
        and not user_known_update_package_registry
      output: >
        Repository files get updated (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid pcmdline=%proc.pcmdline file=%fd.name newpath=%evt.arg.newpath container_id=%container.id image=%container.image.repository)
      priority:
        NOTICE
      tags: [host, container, filesystem, mitre_persistence, T1072]
    
    # Users should overwrite this macro to specify conditions under which a
    # write under the binary dir is ignored. For example, it may be okay to
    # install a binary in the context of a ci/cd build.
    - macro: user_known_write_below_binary_dir_activities
      condition: (never_true)
    
    - rule: Write below binary dir
      desc: an attempt to write to any file below a set of binary directories
      condition: >
        bin_dir and evt.dir = < and open_write
        and not package_mgmt_procs
        and not exe_running_docker_save
        and not python_running_get_pip
        and not python_running_ms_oms
        and not user_known_write_below_binary_dir_activities
      output: >
        File below a known binary directory opened for writing (user=%user.name user_loginuid=%user.loginuid
        command=%proc.cmdline pid=%proc.pid file=%fd.name parent=%proc.pname pcmdline=%proc.pcmdline gparent=%proc.aname[2] container_id=%container.id image=%container.image.repository)
      priority: ERROR
      tags: [host, container, filesystem, mitre_persistence, T1543]
    
    # If you'd like to generally monitor a wider set of directories on top
    # of the ones covered by the rule Write below binary dir, you can use
    # the following rule and lists.
    - list: monitored_directories
      items: [/boot, /lib, /lib64, /usr/lib, /usr/local/lib, /usr/local/sbin, /usr/local/bin, /root/.ssh]
    
    - macro: user_ssh_directory
      condition: (fd.name contains '/.ssh/' and fd.name glob '/home/*/.ssh/*')
    
    - macro: directory_traversal
      condition: (fd.nameraw contains '../' and fd.nameraw glob '*../*../*')
    
    # google_accounts_(daemon)
    - macro: google_accounts_daemon_writing_ssh
      condition: (proc.name=google_accounts and user_ssh_directory)
    
    - macro: cloud_init_writing_ssh
      condition: (proc.name=cloud-init and user_ssh_directory)
    
    - macro: mkinitramfs_writing_boot
      condition: (proc.pname in (mkinitramfs, update-initramf) and fd.directory=/boot)
    
    - macro: monitored_dir
      condition: >
        (fd.directory in (monitored_directories)
         or user_ssh_directory)
        and not mkinitramfs_writing_boot
    
    # Add conditions to this macro (probably in a separate file,
    # overwriting this macro) to allow for specific combinations of
    # programs writing below monitored directories.
    #
    # Its default value is an expression that always is false, which
    # becomes true when the "not ..." in the rule is applied.
    - macro: user_known_write_monitored_dir_conditions
      condition: (never_true)
    
    - rule: Write below monitored dir
      desc: an attempt to write to any file below a set of monitored directories
      condition: >
        evt.dir = < and open_write and monitored_dir
        and not package_mgmt_procs
        and not coreos_write_ssh_dir
        and not exe_running_docker_save
        and not python_running_get_pip
        and not python_running_ms_oms
        and not google_accounts_daemon_writing_ssh
        and not cloud_init_writing_ssh
        and not user_known_write_monitored_dir_conditions
      output: >
        File below a monitored directory opened for writing (user=%user.name user_loginuid=%user.loginuid
        command=%proc.cmdline pid=%proc.pid file=%fd.name parent=%proc.pname pcmdline=%proc.pcmdline gparent=%proc.aname[2] container_id=%container.id image=%container.image.repository)
      priority: ERROR
      tags: [host, container, filesystem, mitre_persistence, T1543]
    
    # ******************************************************************************
    # * "Directory traversal monitored file read" requires FALCO_ENGINE_VERSION 13 *
    # ******************************************************************************
    
    - rule: Directory traversal monitored file read
      desc: >
        Web applications can be vulnerable to directory traversal attacks that allow accessing files outside of the web app's root directory (e.g. Arbitrary File Read bugs).
        System directories like /etc are typically accessed via absolute paths. Access patterns outside of this (here path traversal) can be regarded as suspicious.
        This rule includes failed file open attempts.
      condition: (open_read or open_file_failed) and (etc_dir or user_ssh_directory or fd.name startswith /root/.ssh or fd.name contains "id_rsa") and directory_traversal and not proc.pname in (shell_binaries)
      enabled: true
      output: >
        Read monitored file via directory traversal (username=%user.name useruid=%user.uid user_loginuid=%user.loginuid program=%proc.name exe=%proc.exepath
        command=%proc.cmdline pid=%proc.pid parent=%proc.pname file=%fd.name fileraw=%fd.nameraw parent=%proc.pname
        gparent=%proc.aname[2] container_id=%container.id image=%container.image.repository returncode=%evt.res cwd=%proc.cwd)
      priority: WARNING
      tags: [host, container, filesystem, mitre_discovery, mitre_exfiltration, mitre_credential_access, T1555, T1212, T1020, T1552, T1083]
    
    # The rule below is disabled by default as many system management tools
    # like ansible, etc can read these files/paths. Enable it using this macro.
    - macro: user_known_read_ssh_information_activities
      condition: (never_true)
    
    - rule: Read ssh information
      desc: Any attempt to read files below ssh directories by non-ssh programs
      condition: >
        ((open_read or open_directory) and
         (user_ssh_directory or fd.name startswith /root/.ssh) and
         not user_known_read_ssh_information_activities and
         not proc.name in (ssh_binaries))
      enabled: false
      output: >
        ssh-related file/directory read by non-ssh program (user=%user.name user_loginuid=%user.loginuid
        command=%proc.cmdline pid=%proc.pid file=%fd.name parent=%proc.pname pcmdline=%proc.pcmdline container_id=%container.id image=%container.image.repository)
      priority: ERROR
      tags: [host, container, filesystem, mitre_discovery, T1005]
    
    - list: safe_etc_dirs
      items: [/etc/cassandra, /etc/ssl/certs/java, /etc/logstash, /etc/nginx/conf.d, /etc/container_environment, /etc/hrmconfig, /etc/fluent/configs.d. /etc/alertmanager]
    
    - macro: fluentd_writing_conf_files
      condition: (proc.name=start-fluentd and fd.name in (/etc/fluent/fluent.conf, /etc/td-agent/td-agent.conf))
    
    - macro: qualys_writing_conf_files
      condition: (proc.name=qualys-cloud-ag and fd.name=/etc/qualys/cloud-agent/qagent-log.conf)
    
    - macro: git_writing_nssdb
      condition: (proc.name=git-remote-http and fd.directory=/etc/pki/nssdb)
    
    - macro: plesk_writing_keys
      condition: (proc.name in (plesk_binaries) and fd.name startswith /etc/sw/keys)
    
    - macro: plesk_install_writing_apache_conf
      condition: (proc.cmdline startswith "bash -hB /usr/lib/plesk-9.0/services/webserver.apache configure"
                  and fd.name="/etc/apache2/apache2.conf.tmp")
    
    - macro: plesk_running_mktemp
      condition: (proc.name=mktemp and proc.aname[3] in (plesk_binaries))
    
    - macro: networkmanager_writing_resolv_conf
      condition: proc.aname[2]=nm-dispatcher and fd.name=/etc/resolv.conf
    
    - macro: add_shell_writing_shells_tmp
      condition: (proc.name=add-shell and fd.name=/etc/shells.tmp)
    
    - macro: duply_writing_exclude_files
      condition: (proc.name=touch and proc.pcmdline startswith "bash /usr/bin/duply" and fd.name startswith "/etc/duply")
    
    - macro: xmlcatalog_writing_files
      condition: (proc.name=update-xmlcatal and fd.directory=/etc/xml)
    
    - macro: datadog_writing_conf
      condition: ((proc.cmdline startswith "python /opt/datadog-agent" or
                   proc.cmdline startswith "entrypoint.sh /entrypoint.sh datadog start" or
                   proc.cmdline startswith "agent.py /opt/datadog-agent")
                  and fd.name startswith "/etc/dd-agent")
    
    - macro: rancher_writing_conf
      condition: ((proc.name in (healthcheck, lb-controller, rancher-dns)) and
                  (container.image.repository contains "rancher/healthcheck" or
                   container.image.repository contains "rancher/lb-service-haproxy" or
                   container.image.repository contains "rancher/dns") and
                  (fd.name startswith "/etc/haproxy" or fd.name startswith "/etc/rancher-dns"))
    
    - macro: rancher_writing_root
      condition: (proc.name=rancher-metadat and
                  (container.image.repository contains "rancher/metadata" or container.image.repository contains "rancher/lb-service-haproxy") and
                  fd.name startswith "/answers.json")
    
    - macro: checkpoint_writing_state
      condition: (proc.name=checkpoint and
                  container.image.repository contains "coreos/pod-checkpointer" and
                  fd.name startswith "/etc/kubernetes")
    
    - macro: jboss_in_container_writing_passwd
      condition: >
        ((proc.cmdline="run-java.sh /opt/jboss/container/java/run/run-java.sh"
          or proc.cmdline="run-java.sh /opt/run-java/run-java.sh")
         and container
         and fd.name=/etc/passwd)
    
    - macro: curl_writing_pki_db
      condition: (proc.name=curl and fd.directory=/etc/pki/nssdb)
    
    - macro: haproxy_writing_conf
      condition: ((proc.name in (update-haproxy-,haproxy_reload.) or proc.pname in (update-haproxy-,haproxy_reload,haproxy_reload.))
                   and (fd.name=/etc/openvpn/client.map or fd.name startswith /etc/haproxy))
    
    - macro: java_writing_conf
      condition: (proc.name=java and fd.name=/etc/.java/.systemPrefs/.system.lock)
    
    - macro: rabbitmq_writing_conf
      condition: (proc.name=rabbitmq-server and fd.directory=/etc/rabbitmq)
    
    - macro: rook_writing_conf
      condition: (proc.name=toolbox.sh and container.image.repository=rook/toolbox
                  and fd.directory=/etc/ceph)
    
    - macro: httpd_writing_conf_logs
      condition: (proc.name=httpd and fd.name startswith /etc/httpd/)
    
    - macro: mysql_writing_conf
      condition: >
        ((proc.name in (start-mysql.sh, run-mysqld) or proc.pname=start-mysql.sh) and
         (fd.name startswith /etc/mysql or fd.directory=/etc/my.cnf.d))
    
    - macro: redis_writing_conf
      condition: >
        (proc.name in (run-redis, redis-launcher.) and (fd.name=/etc/redis.conf or fd.name startswith /etc/redis))
    
    - macro: openvpn_writing_conf
      condition: (proc.name in (openvpn,openvpn-entrypo) and fd.name startswith /etc/openvpn)
    
    - macro: php_handlers_writing_conf
      condition: (proc.name=php_handlers_co and fd.name=/etc/psa/php_versions.json)
    
    - macro: sed_writing_temp_file
      condition: >
        ((proc.aname[3]=cron_start.sh and fd.name startswith /etc/security/sed) or
         (proc.name=sed and (fd.name startswith /etc/apt/sources.list.d/sed or
                             fd.name startswith /etc/apt/sed or
                             fd.name startswith /etc/apt/apt.conf.d/sed)))
    
    - macro: cron_start_writing_pam_env
      condition: (proc.cmdline="bash /usr/sbin/start-cron" and fd.name=/etc/security/pam_env.conf)
    
    # In some cases dpkg-reconfigur runs commands that modify /etc. Not
    # putting the full set of package management programs yet.
    - macro: dpkg_scripting
      condition: (proc.aname[2] in (dpkg-reconfigur, dpkg-preconfigu))
    
    - macro: ufw_writing_conf
      condition: (proc.name=ufw and fd.directory=/etc/ufw)
    
    - macro: calico_writing_conf
      condition: >
        (((proc.name = calico-node) or
          (container.image.repository=gcr.io/projectcalico-org/node and proc.name in (start_runit, cp)) or
          (container.image.repository=gcr.io/projectcalico-org/cni and proc.name=sed))
         and fd.name startswith /etc/calico)
    
    - macro: prometheus_conf_writing_conf
      condition: (proc.name=prometheus-conf and fd.name startswith /etc/prometheus/config_out)
    
    - macro: openshift_writing_conf
      condition: (proc.name=oc and fd.name startswith /etc/origin/node)
    
    - macro: keepalived_writing_conf
      condition: (proc.name in (keepalived, kube-keepalived) and fd.name=/etc/keepalived/keepalived.conf)
    
    - macro: etcd_manager_updating_dns
      condition: (container and proc.name=etcd-manager and fd.name=/etc/hosts)
    
    - macro: automount_using_mtab
      condition: (proc.pname = automount and fd.name startswith /etc/mtab)
    
    - macro: mcafee_writing_cma_d
      condition: (proc.name=macompatsvc and fd.directory=/etc/cma.d)
    
    - macro: avinetworks_supervisor_writing_ssh
      condition: >
        (proc.cmdline="se_supervisor.p /opt/avi/scripts/se_supervisor.py -d" and
          (fd.name startswith /etc/ssh/known_host_ or
           fd.name startswith /etc/ssh/ssh_monitor_config_ or
           fd.name startswith /etc/ssh/ssh_config_))
    
    - macro: multipath_writing_conf
      condition: (proc.name = multipath and fd.name startswith /etc/multipath/)
    
    # Add conditions to this macro (probably in a separate file,
    # overwriting this macro) to allow for specific combinations of
    # programs writing below specific directories below
    # /etc. fluentd_writing_conf_files is a good example to follow, as it
    # specifies both the program doing the writing as well as the specific
    # files it is allowed to modify.
    #
    # In this file, it just takes one of the programs in the base macro
    # and repeats it.
    
    - macro: user_known_write_etc_conditions
      condition: proc.name=confd
    
    # This is a placeholder for user to extend the whitelist for write below etc rule
    - macro: user_known_write_below_etc_activities
      condition: (never_true)
    
    - macro: calico_node
      condition: (container.image.repository endswith calico/node and proc.name=calico-node)
    
    - macro: write_etc_common
      condition: >
        etc_dir and evt.dir = < and open_write
        and proc_name_exists
        and not proc.name in (passwd_binaries, shadowutils_binaries, sysdigcloud_binaries,
                              package_mgmt_binaries, ssl_mgmt_binaries, dhcp_binaries,
                              dev_creation_binaries, shell_mgmt_binaries,
                              mail_config_binaries,
                              sshkit_script_binaries,
                              ldconfig.real, ldconfig, confd, gpg, insserv,
                              apparmor_parser, update-mime, tzdata.config, tzdata.postinst,
                              systemd, systemd-machine, systemd-sysuser,
                              debconf-show, rollerd, bind9.postinst, sv,
                              gen_resolvconf., update-ca-certi, certbot, runsv,
                              qualys-cloud-ag, locales.postins, nomachine_binaries,
                              adclient, certutil, crlutil, pam-auth-update, parallels_insta,
                              openshift-launc, update-rc.d, puppet, falcoctl)
        and not (container and proc.cmdline in ("cp /run/secrets/kubernetes.io/serviceaccount/ca.crt /etc/pki/ca-trust/source/anchors/openshift-ca.crt"))
        and not proc.pname in (sysdigcloud_binaries, mail_config_binaries, hddtemp.postins, sshkit_script_binaries, locales.postins, deb_binaries, dhcp_binaries)
        and not fd.name pmatch (safe_etc_dirs)
        and not fd.name in (/etc/container_environment.sh, /etc/container_environment.json, /etc/motd, /etc/motd.svc)
        and not sed_temporary_file
        and not exe_running_docker_save
        and not ansible_running_python
        and not python_running_denyhosts
        and not fluentd_writing_conf_files
        and not user_known_write_etc_conditions
        and not run_by_centrify
        and not run_by_adclient
        and not qualys_writing_conf_files
        and not git_writing_nssdb
        and not plesk_writing_keys
        and not plesk_install_writing_apache_conf
        and not plesk_running_mktemp
        and not networkmanager_writing_resolv_conf
        and not run_by_chef
        and not add_shell_writing_shells_tmp
        and not duply_writing_exclude_files
        and not xmlcatalog_writing_files
        and not parent_supervise_running_multilog
        and not supervise_writing_status
        and not pki_realm_writing_realms
        and not htpasswd_writing_passwd
        and not lvprogs_writing_conf
        and not ovsdb_writing_openvswitch
        and not datadog_writing_conf
        and not curl_writing_pki_db
        and not haproxy_writing_conf
        and not java_writing_conf
        and not dpkg_scripting
        and not parent_ucf_writing_conf
        and not rabbitmq_writing_conf
        and not rook_writing_conf
        and not php_handlers_writing_conf
        and not sed_writing_temp_file
        and not cron_start_writing_pam_env
        and not httpd_writing_conf_logs
        and not mysql_writing_conf
        and not openvpn_writing_conf
        and not consul_template_writing_conf
        and not countly_writing_nginx_conf
        and not ms_oms_writing_conf
        and not ms_scx_writing_conf
        and not azure_scripts_writing_conf
        and not azure_networkwatcher_writing_conf
        and not couchdb_writing_conf
        and not update_texmf_writing_conf
        and not slapadd_writing_conf
        and not symantec_writing_conf
        and not liveupdate_writing_conf
        and not sosreport_writing_files
        and not selinux_writing_conf
        and not veritas_writing_config
        and not nginx_writing_conf
        and not nginx_writing_certs
        and not chef_client_writing_conf
        and not centrify_writing_krb
        and not sssd_writing_krb
        and not cockpit_writing_conf
        and not ipsec_writing_conf
        and not httpd_writing_ssl_conf
        and not userhelper_writing_etc_security
        and not pkgmgmt_progs_writing_pki
        and not update_ca_trust_writing_pki
        and not brandbot_writing_os_release
        and not redis_writing_conf
        and not openldap_writing_conf
        and not ucpagent_writing_conf
        and not iscsi_writing_conf
        and not istio_writing_conf
        and not ufw_writing_conf
        and not calico_writing_conf
        and not calico_writing_envvars
        and not prometheus_conf_writing_conf
        and not openshift_writing_conf
        and not keepalived_writing_conf
        and not rancher_writing_conf
        and not checkpoint_writing_state
        and not jboss_in_container_writing_passwd
        and not etcd_manager_updating_dns
        and not user_known_write_below_etc_activities
        and not automount_using_mtab
        and not mcafee_writing_cma_d
        and not avinetworks_supervisor_writing_ssh
        and not multipath_writing_conf
        and not calico_node
    
    - rule: Write below etc
      desc: an attempt to write to any file below /etc
      condition: write_etc_common
      output: "File below /etc opened for writing (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid parent=%proc.pname pcmdline=%proc.pcmdline file=%fd.name program=%proc.name gparent=%proc.aname[2] ggparent=%proc.aname[3] gggparent=%proc.aname[4] container_id=%container.id image=%container.image.repository)"
      priority: ERROR
      tags: [host, container, filesystem, mitre_persistence, T1098]
    
    - list: known_root_files
      items: [/root/.monit.state, /root/.auth_tokens, /root/.bash_history, /root/.ash_history, /root/.aws/credentials,
              /root/.viminfo.tmp, /root/.lesshst, /root/.bzr.log, /root/.gitconfig.lock, /root/.babel.json, /root/.localstack,
              /root/.node_repl_history, /root/.mongorc.js, /root/.dbshell, /root/.augeas/history, /root/.rnd, /root/.wget-hsts, /health, /exec.fifo]
    
    - list: known_root_directories
      items: [/root/.oracle_jre_usage, /root/.ssh, /root/.subversion, /root/.nami]
    
    - macro: known_root_conditions
      condition: (fd.name startswith /root/orcexec.
                  or fd.name startswith /root/.m2
                  or fd.name startswith /root/.npm
                  or fd.name startswith /root/.pki
                  or fd.name startswith /root/.ivy2
                  or fd.name startswith /root/.config/Cypress
                  or fd.name startswith /root/.config/pulse
                  or fd.name startswith /root/.config/configstore
                  or fd.name startswith /root/jenkins/workspace
                  or fd.name startswith /root/.jenkins
                  or fd.name startswith /root/.cache
                  or fd.name startswith /root/.sbt
                  or fd.name startswith /root/.java
                  or fd.name startswith /root/.glide
                  or fd.name startswith /root/.sonar
                  or fd.name startswith /root/.v8flag
                  or fd.name startswith /root/infaagent
                  or fd.name startswith /root/.local/lib/python
                  or fd.name startswith /root/.pm2
                  or fd.name startswith /root/.gnupg
                  or fd.name startswith /root/.pgpass
                  or fd.name startswith /root/.theano
                  or fd.name startswith /root/.gradle
                  or fd.name startswith /root/.android
                  or fd.name startswith /root/.ansible
                  or fd.name startswith /root/.crashlytics
                  or fd.name startswith /root/.dbus
                  or fd.name startswith /root/.composer
                  or fd.name startswith /root/.gconf
                  or fd.name startswith /root/.nv
                  or fd.name startswith /root/.local/share/jupyter
                  or fd.name startswith /root/oradiag_root
                  or fd.name startswith /root/workspace
                  or fd.name startswith /root/jvm
                  or fd.name startswith /root/.node-gyp)
    
    # Add conditions to this macro (probably in a separate file,
    # overwriting this macro) to allow for specific combinations of
    # programs writing below specific directories below
    # / or /root.
    #
    # In this file, it just takes one of the condition in the base macro
    # and repeats it.
    - macro: user_known_write_root_conditions
      condition: fd.name=/root/.bash_history
    
    # This is a placeholder for user to extend the whitelist for write below root rule
    - macro: user_known_write_below_root_activities
      condition: (never_true)
    
    - macro: runc_writing_exec_fifo
      condition: (proc.cmdline="runc:[1:CHILD] init" and fd.name=/exec.fifo)
    
    - macro: runc_writing_var_lib_docker
      condition: (proc.cmdline="runc:[1:CHILD] init" and evt.arg.filename startswith /var/lib/docker)
    
    - macro: mysqlsh_writing_state
      condition: (proc.name=mysqlsh and fd.directory=/root/.mysqlsh)
    
    - rule: Write below root
      desc: an attempt to write to any file directly below / or /root
      condition: >
        root_dir and evt.dir = < and open_write
        and proc_name_exists
        and not fd.name in (known_root_files)
        and not fd.directory pmatch (known_root_directories)
        and not exe_running_docker_save
        and not gugent_writing_guestagent_log
        and not dse_writing_tmp
        and not zap_writing_state
        and not airflow_writing_state
        and not rpm_writing_root_rpmdb
        and not maven_writing_groovy
        and not chef_writing_conf
        and not kubectl_writing_state
        and not cassandra_writing_state
        and not galley_writing_state
        and not calico_writing_state
        and not rancher_writing_root
        and not runc_writing_exec_fifo
        and not mysqlsh_writing_state
        and not known_root_conditions
        and not user_known_write_root_conditions
        and not user_known_write_below_root_activities
      output: "File below / or /root opened for writing (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid parent=%proc.pname file=%fd.name program=%proc.name container_id=%container.id image=%container.image.repository)"
      priority: ERROR
      tags: [host, container, filesystem, mitre_persistence, TA0003]
    
    - macro: cmp_cp_by_passwd
      condition: proc.name in (cmp, cp) and proc.pname in (passwd, run-parts)
    
    - macro: user_known_read_sensitive_files_activities
      condition: (never_true)
    
    - rule: Read sensitive file trusted after startup
      desc: >
        an attempt to read any sensitive file (e.g. files containing user/password/authentication
        information) by a trusted program after startup. Trusted programs might read these files
        at startup to load initial state, but not afterwards.
      condition: sensitive_files and open_read and server_procs and not proc_is_new and proc.name!="sshd" and not user_known_read_sensitive_files_activities
      output: >
        Sensitive file opened for reading by trusted program after startup (user=%user.name user_loginuid=%user.loginuid
        command=%proc.cmdline pid=%proc.pid parent=%proc.pname file=%fd.name parent=%proc.pname gparent=%proc.aname[2] container_id=%container.id image=%container.image.repository)
      priority: WARNING
      tags: [host, container, filesystem, mitre_credential_access, T1555, T1212, T1020, T1552, T1083]
    
    - list: read_sensitive_file_binaries
      items: [
        iptables, ps, lsb_release, check-new-relea, dumpe2fs, accounts-daemon, sshd,
        vsftpd, systemd, mysql_install_d, psql, screen, debconf-show, sa-update,
        pam-auth-update, pam-config, /usr/sbin/spamd, polkit-agent-he, lsattr, file, sosreport,
        scxcimservera, adclient, rtvscand, cockpit-session, userhelper, ossec-syscheckd
        ]
    
    # Add conditions to this macro (probably in a separate file,
    # overwriting this macro) to allow for specific combinations of
    # programs accessing sensitive files.
    # fluentd_writing_conf_files is a good example to follow, as it
    # specifies both the program doing the writing as well as the specific
    # files it is allowed to modify.
    #
    # In this file, it just takes one of the macros in the base rule
    # and repeats it.
    
    - macro: user_read_sensitive_file_conditions
      condition: cmp_cp_by_passwd
    
    - list: read_sensitive_file_images
      items: []
    
    - macro: user_read_sensitive_file_containers
      condition: (container and container.image.repository in (read_sensitive_file_images))
    
    # This macro detects man-db postinst, see https://salsa.debian.org/debian/man-db/-/blob/master/debian/postinst
    # The rule "Read sensitive file untrusted" use this macro to avoid FPs.
    - macro: mandb_postinst
      condition: >
        (proc.name=perl and proc.args startswith "-e" and
        proc.args contains "@pwd = getpwnam(" and
        proc.args contains "exec " and
        proc.args contains "/usr/bin/mandb")
    
    - rule: Read sensitive file untrusted
      desc: >
        an attempt to read any sensitive file (e.g. files containing user/password/authentication
        information). Exceptions are made for known trusted programs.
      condition: >
        sensitive_files and open_read
        and proc_name_exists
        and not proc.name in (user_mgmt_binaries, userexec_binaries, package_mgmt_binaries,
         cron_binaries, read_sensitive_file_binaries, shell_binaries, hids_binaries,
         vpn_binaries, mail_config_binaries, nomachine_binaries, sshkit_script_binaries,
         in.proftpd, mandb, salt-call, salt-minion, postgres_mgmt_binaries,
         google_oslogin_
         )
        and not cmp_cp_by_passwd
        and not ansible_running_python
        and not run_by_qualys
        and not run_by_chef
        and not run_by_google_accounts_daemon
        and not user_read_sensitive_file_conditions
        and not mandb_postinst
        and not perl_running_plesk
        and not perl_running_updmap
        and not veritas_driver_script
        and not perl_running_centrifydc
        and not runuser_reading_pam
        and not linux_bench_reading_etc_shadow
        and not user_known_read_sensitive_files_activities
        and not user_read_sensitive_file_containers
      output: >
        Sensitive file opened for reading by non-trusted program (user=%user.name user_loginuid=%user.loginuid program=%proc.name
        command=%proc.cmdline pid=%proc.pid file=%fd.name parent=%proc.pname gparent=%proc.aname[2] ggparent=%proc.aname[3] gggparent=%proc.aname[4] container_id=%container.id image=%container.image.repository)
      priority: WARNING
      tags: [host, container, filesystem, mitre_credential_access, mitre_discovery, T1555, T1212, T1020, T1552, T1083]
    
    - macro: amazon_linux_running_python_yum
      condition: >
        (proc.name = python and
         proc.pcmdline = "python -m amazon_linux_extras system_motd" and
         proc.cmdline startswith "python -c import yum;")
    
    - macro: user_known_write_rpm_database_activities
      condition: (never_true)
    
    # Only let rpm-related programs write to the rpm database
    - rule: Write below rpm database
      desc: an attempt to write to the rpm database by any non-rpm related program
      condition: >
        fd.name startswith /var/lib/rpm and open_write
        and not rpm_procs
        and not ansible_running_python
        and not python_running_chef
        and not exe_running_docker_save
        and not amazon_linux_running_python_yum
        and not user_known_write_rpm_database_activities
      output: "Rpm database opened for writing by a non-rpm program (command=%proc.cmdline pid=%proc.pid file=%fd.name parent=%proc.pname pcmdline=%proc.pcmdline container_id=%container.id image=%container.image.repository)"
      priority: ERROR
      tags: [host, container, filesystem, software_mgmt, mitre_persistence, T1072]
    
    - macro: postgres_running_wal_e
      condition: (proc.pname=postgres and (proc.cmdline startswith "sh -c envdir /etc/wal-e.d/env /usr/local/bin/wal-e" or proc.cmdline startswith "sh -c envdir \"/run/etc/wal-e.d/env\" wal-g wal-push"))
    
    - macro: redis_running_prepost_scripts
      condition: (proc.aname[2]=redis-server and (proc.cmdline contains "redis-server.post-up.d" or proc.cmdline contains "redis-server.pre-up.d"))
    
    - macro: rabbitmq_running_scripts
      condition: >
        (proc.pname=beam.smp and
        (proc.cmdline startswith "sh -c exec ps" or
         proc.cmdline startswith "sh -c exec inet_gethost" or
         proc.cmdline= "sh -s unix:cmd" or
         proc.cmdline= "sh -c exec /bin/sh -s unix:cmd 2>&1"))
    
    - macro: rabbitmqctl_running_scripts
      condition: (proc.aname[2]=rabbitmqctl and proc.cmdline startswith "sh -c ")
    
    - macro: run_by_appdynamics
      condition: (proc.pname=java and proc.pcmdline startswith "java -jar -Dappdynamics")
    
    - macro: user_known_db_spawned_processes
      condition: (never_true)
    
    - rule: DB program spawned process
      desc: >
        a database-server related program spawned a new process other than itself.
        This shouldn\'t occur and is a follow on from some SQL injection attacks.
      condition: >
        proc.pname in (db_server_binaries)
        and spawned_process
        and not proc.name in (db_server_binaries)
        and not postgres_running_wal_e
        and not user_known_db_spawned_processes
      output: >
        Database-related program spawned process other than itself (user=%user.name user_loginuid=%user.loginuid
        program=%proc.cmdline pid=%proc.pid parent=%proc.pname container_id=%container.id image=%container.image.repository exe_flags=%evt.arg.flags)
      priority: NOTICE
      tags: [host, container, process, database, mitre_execution, T1190]
    
    - macro: user_known_modify_bin_dir_activities
      condition: (never_true)
    
    - rule: Modify binary dirs
      desc: an attempt to modify any file below a set of binary directories.
      condition: bin_dir_rename and modify and not package_mgmt_procs and not exe_running_docker_save and not user_known_modify_bin_dir_activities
      output: >
        File below known binary directory renamed/removed (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid
        pcmdline=%proc.pcmdline operation=%evt.type file=%fd.name %evt.args container_id=%container.id image=%container.image.repository)
      priority: ERROR
      tags: [host, container, filesystem, mitre_persistence, T1222.002]
    
    - macro: user_known_mkdir_bin_dir_activities
      condition: (never_true)
    
    - rule: Mkdir binary dirs
      desc: an attempt to create a directory below a set of binary directories.
      condition: >
        mkdir
        and bin_dir_mkdir
        and not package_mgmt_procs
        and not user_known_mkdir_bin_dir_activities
        and not exe_running_docker_save
      output: >
        Directory below known binary directory created (user=%user.name user_loginuid=%user.loginuid
        command=%proc.cmdline pid=%proc.pid directory=%evt.arg.path container_id=%container.id image=%container.image.repository)
      priority: ERROR
      tags: [host, container, filesystem, mitre_persistence, T1222.002]
    
    # This list allows for easy additions to the set of commands allowed
    # to change thread namespace without having to copy and override the
    # entire change thread namespace rule.
    - list: user_known_change_thread_namespace_binaries
      items: [crio, multus]
    
    - macro: user_known_change_thread_namespace_activities
      condition: (never_true)
    
    - list: network_plugin_binaries
      items: [aws-cni, azure-vnet]
    
    - macro: weaveworks_scope
      condition: (container.image.repository endswith weaveworks/scope and proc.name=scope)
    
    - rule: Change thread namespace
      desc: >
        an attempt to change a program/thread\'s namespace (commonly done
        as a part of creating a container) by calling setns.
      condition: >
        evt.type=setns and evt.dir=<
        and proc_name_exists
        and not (container.id=host and proc.name in (docker_binaries, k8s_binaries, lxd_binaries, nsenter))
        and not proc.name in (sysdigcloud_binaries, sysdig, calico, oci-umount, cilium-cni, network_plugin_binaries)
        and not proc.name in (user_known_change_thread_namespace_binaries)
        and not proc.name startswith "runc"
        and not proc.cmdline startswith "containerd"
        and not proc.pname in (sysdigcloud_binaries, hyperkube, kubelet, protokube, dockerd, tini, aws)
        and not java_running_sdjagent
        and not kubelet_running_loopback
        and not rancher_agent
        and not rancher_network_manager
        and not calico_node
        and not weaveworks_scope
        and not user_known_change_thread_namespace_activities
      enabled: false
      output: >
        Namespace change (setns) by unexpected program (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid
        parent=%proc.pname %container.info container_id=%container.id image=%container.image.repository:%container.image.tag)
      priority: NOTICE
      tags: [host, container, process, mitre_privilege_escalation, mitre_lateral_movement, T1611]
    
    # The binaries in this list and their descendents are *not* allowed
    # spawn shells. This includes the binaries spawning shells directly as
    # well as indirectly. For example, apache -> php/perl for
    # mod_{php,perl} -> some shell is also not allowed, because the shell
    # has apache as an ancestor.
    
    - list: protected_shell_spawning_binaries
      items: [
        http_server_binaries, db_server_binaries, nosql_server_binaries, mail_binaries,
        fluentd, flanneld, splunkd, consul, smbd, runsv, PM2
        ]
    
    - macro: parent_java_running_zookeeper
      condition: (proc.pname=java and proc.pcmdline contains org.apache.zookeeper.server)
    
    - macro: parent_java_running_kafka
      condition: (proc.pname=java and proc.pcmdline contains kafka.Kafka)
    
    - macro: parent_java_running_elasticsearch
      condition: (proc.pname=java and proc.pcmdline contains org.elasticsearch.bootstrap.Elasticsearch)
    
    - macro: parent_java_running_activemq
      condition: (proc.pname=java and proc.pcmdline contains activemq.jar)
    
    - macro: parent_java_running_cassandra
      condition: (proc.pname=java and (proc.pcmdline contains "-Dcassandra.config.loader" or proc.pcmdline contains org.apache.cassandra.service.CassandraDaemon))
    
    - macro: parent_java_running_jboss_wildfly
      condition: (proc.pname=java and proc.pcmdline contains org.jboss)
    
    - macro: parent_java_running_glassfish
      condition: (proc.pname=java and proc.pcmdline contains com.sun.enterprise.glassfish)
    
    - macro: parent_java_running_hadoop
      condition: (proc.pname=java and proc.pcmdline contains org.apache.hadoop)
    
    - macro: parent_java_running_datastax
      condition: (proc.pname=java and proc.pcmdline contains com.datastax)
    
    - macro: nginx_starting_nginx
      condition: (proc.pname=nginx and proc.cmdline contains "/usr/sbin/nginx -c /etc/nginx/nginx.conf")
    
    - macro: nginx_running_aws_s3_cp
      condition: (proc.pname=nginx and proc.cmdline startswith "sh -c /usr/local/bin/aws s3 cp")
    
    - macro: consul_running_net_scripts
      condition: (proc.pname=consul and (proc.cmdline startswith "sh -c curl" or proc.cmdline startswith "sh -c nc"))
    
    - macro: consul_running_alert_checks
      condition: (proc.pname=consul and proc.cmdline startswith "sh -c /bin/consul-alerts")
    
    - macro: serf_script
      condition: (proc.cmdline startswith "sh -c serf")
    
    - macro: check_process_status
      condition: (proc.cmdline startswith "sh -c kill -0 ")
    
    # In some cases, you may want to consider node processes run directly
    # in containers as protected shell spawners. Examples include using
    # pm2-docker or pm2 start some-app.js --no-daemon-mode as the direct
    # entrypoint of the container, and when the node app is a long-lived
    # server using something like express.
    #
    # However, there are other uses of node related to build pipelines for
    # which node is not really a server but instead a general scripting
    # tool. In these cases, shells are very likely and in these cases you
    # don't want to consider node processes protected shell spawners.
    #
    # We have to choose one of these cases, so we consider node processes
    # as unprotected by default. If you want to consider any node process
    # run in a container as a protected shell spawner, override the below
    # macro to remove the "never_true" clause, which allows it to take effect.
    - macro: possibly_node_in_container
      condition: (never_true and (proc.pname=node and proc.aname[3]=docker-containe))
    
    # Similarly, you may want to consider any shell spawned by apache
    # tomcat as suspect. The famous apache struts attack (CVE-2017-5638)
    # could be exploited to do things like spawn shells.
    #
    # However, many applications *do* use tomcat to run arbitrary shells,
    # as a part of build pipelines, etc.
    #
    # Like for node, we make this case opt-in.
    - macro: possibly_parent_java_running_tomcat
      condition: (never_true and proc.pname=java and proc.pcmdline contains org.apache.catalina.startup.Bootstrap)
    
    - macro: protected_shell_spawner
      condition: >
        (proc.aname in (protected_shell_spawning_binaries)
        or parent_java_running_zookeeper
        or parent_java_running_kafka
        or parent_java_running_elasticsearch
        or parent_java_running_activemq
        or parent_java_running_cassandra
        or parent_java_running_jboss_wildfly
        or parent_java_running_glassfish
        or parent_java_running_hadoop
        or parent_java_running_datastax
        or possibly_parent_java_running_tomcat
        or possibly_node_in_container)
    
    - list: mesos_shell_binaries
      items: [mesos-docker-ex, mesos-slave, mesos-health-ch]
    
    # Note that runsv is both in protected_shell_spawner and the
    # exclusions by pname. This means that runsv can itself spawn shells
    # (the ./run and ./finish scripts), but the processes runsv can not
    # spawn shells.
    - rule: Run shell untrusted
      desc: an attempt to spawn a shell below a non-shell application. Specific applications are monitored.
      condition: >
        spawned_process
        and shell_procs
        and proc.pname exists
        and protected_shell_spawner
        and not proc.pname in (shell_binaries, gitlab_binaries, cron_binaries, user_known_shell_spawn_binaries,
                               needrestart_binaries,
                               mesos_shell_binaries,
                               erl_child_setup, exechealthz,
                               PM2, PassengerWatchd, c_rehash, svlogd, logrotate, hhvm, serf,
                               lb-controller, nvidia-installe, runsv, statsite, erlexec, calico-node,
                               "puma reactor")
        and not proc.cmdline in (known_shell_spawn_cmdlines)
        and not proc.aname in (unicorn_launche)
        and not consul_running_net_scripts
        and not consul_running_alert_checks
        and not nginx_starting_nginx
        and not nginx_running_aws_s3_cp
        and not run_by_package_mgmt_binaries
        and not serf_script
        and not check_process_status
        and not run_by_foreman
        and not python_mesos_marathon_scripting
        and not splunk_running_forwarder
        and not postgres_running_wal_e
        and not redis_running_prepost_scripts
        and not rabbitmq_running_scripts
        and not rabbitmqctl_running_scripts
        and not run_by_appdynamics
        and not user_shell_container_exclusions
      output: >
        Shell spawned by untrusted binary (user=%user.name user_loginuid=%user.loginuid shell=%proc.name parent=%proc.pname
        cmdline=%proc.cmdline pid=%proc.pid pcmdline=%proc.pcmdline gparent=%proc.aname[2] ggparent=%proc.aname[3]
        aname[4]=%proc.aname[4] aname[5]=%proc.aname[5] aname[6]=%proc.aname[6] aname[7]=%proc.aname[7] container_id=%container.id image=%container.image.repository exe_flags=%evt.arg.flags)
      priority: DEBUG
      tags: [host, container, process, shell, mitre_execution, T1059.004]
    
    - macro: allowed_openshift_registry_root
      condition: >
        (container.image.repository startswith openshift3/ or
         container.image.repository startswith registry.redhat.io/openshift3/ or
         container.image.repository startswith registry.access.redhat.com/openshift3/)
    
    # Source: https://docs.openshift.com/enterprise/3.2/install_config/install/disconnected_install.html
    - macro: openshift_image
      condition: >
        (allowed_openshift_registry_root and
          (container.image.repository endswith /logging-deployment or
           container.image.repository endswith /logging-elasticsearch or
           container.image.repository endswith /logging-kibana or
           container.image.repository endswith /logging-fluentd or
           container.image.repository endswith /logging-auth-proxy or
           container.image.repository endswith /metrics-deployer or
           container.image.repository endswith /metrics-hawkular-metrics or
           container.image.repository endswith /metrics-cassandra or
           container.image.repository endswith /metrics-heapster or
           container.image.repository endswith /ose-haproxy-router or
           container.image.repository endswith /ose-deployer or
           container.image.repository endswith /ose-sti-builder or
           container.image.repository endswith /ose-docker-builder or
           container.image.repository endswith /ose-pod or
           container.image.repository endswith /ose-node or
           container.image.repository endswith /ose-docker-registry or
           container.image.repository endswith /prometheus-node-exporter or
           container.image.repository endswith /image-inspector))
    
    - list: redhat_io_images_privileged
      items: [registry.redhat.io/openshift-logging/fluentd-rhel8, registry.redhat.io/openshift4/ose-csi-node-driver-registrar, registry.redhat.io/openshift4/ose-kubernetes-nmstate-handler-rhel8, registry.redhat.io/openshift4/ose-local-storage-diskmaker]
    
    - macro: redhat_image
      condition: >
        (container.image.repository in (redhat_io_images_privileged))
    
    # https://docs.aws.amazon.com/eks/latest/userguide/add-ons-images.html
    #  official AWS EKS registry list. AWS has different ECR repo per region
    - macro: allowed_aws_ecr_registry_root_for_eks
      condition: >
        (container.image.repository startswith "602401143452.dkr.ecr" or
         container.image.repository startswith "877085696533.dkr.ecr" or
         container.image.repository startswith "800184023465.dkr.ecr" or
         container.image.repository startswith "918309763551.dkr.ecr" or
         container.image.repository startswith "961992271922.dkr.ecr" or
         container.image.repository startswith "590381155156.dkr.ecr" or
         container.image.repository startswith "558608220178.dkr.ecr" or
         container.image.repository startswith "151742754352.dkr.ecr" or
         container.image.repository startswith "013241004608.dkr.ecr")
    
    
    - macro: aws_eks_core_images
      condition: >
        (allowed_aws_ecr_registry_root_for_eks and
        (container.image.repository endswith ".amazonaws.com/amazon-k8s-cni" or
         container.image.repository endswith ".amazonaws.com/eks/kube-proxy"))
    
    
    - macro: aws_eks_image_sensitive_mount
      condition: >
        (allowed_aws_ecr_registry_root_for_eks and container.image.repository endswith ".amazonaws.com/amazon-k8s-cni")
    
    # These images are allowed both to run with --privileged and to mount
    # sensitive paths from the host filesystem.
    #
    # NOTE: This list is only provided for backwards compatibility with
    # older local falco rules files that may have been appending to
    # trusted_images. To make customizations, it's better to add images to
    # either privileged_images or falco_sensitive_mount_images.
    - list: trusted_images
      items: []
    
    # Add conditions to this macro (probably in a separate file,
    # overwriting this macro) to specify additional containers that are
    # trusted and therefore allowed to run privileged *and* with sensitive
    # mounts.
    #
    # Like trusted_images, this is deprecated in favor of
    # user_privileged_containers and user_sensitive_mount_containers and
    # is only provided for backwards compatibility.
    #
    # In this file, it just takes one of the images in trusted_containers
    # and repeats it.
    - macro: user_trusted_containers
      condition: (never_true)
    
    - list: sematext_images
      items: [docker.io/sematext/sematext-agent-docker, docker.io/sematext/agent, docker.io/sematext/logagent,
              registry.access.redhat.com/sematext/sematext-agent-docker,
              registry.access.redhat.com/sematext/agent,
              registry.access.redhat.com/sematext/logagent]
    
    # Falco containers
    - list: falco_containers
      items:
        - falcosecurity/falco
        - docker.io/falcosecurity/falco
        - public.ecr.aws/falcosecurity/falco
    
    # Falco no driver containers
    - list: falco_no_driver_containers
      items:
        - falcosecurity/falco-no-driver
        - docker.io/falcosecurity/falco-no-driver
        - public.ecr.aws/falcosecurity/falco-no-driver
    
    # These container images are allowed to run with --privileged and full set of capabilities
    # TODO: Remove k8s.gcr.io reference after 01/Dec/2023
    - list: falco_privileged_images
      items: [
        falco_containers,
        docker.io/calico/node,
        calico/node,
        docker.io/cloudnativelabs/kube-router,
        docker.io/docker/ucp-agent,
        docker.io/mesosphere/mesos-slave,
        docker.io/rook/toolbox,
        docker.io/sysdig/sysdig,
        gcr.io/google_containers/kube-proxy,
        gcr.io/google-containers/startup-script,
        gcr.io/projectcalico-org/node,
        gke.gcr.io/kube-proxy,
        gke.gcr.io/gke-metadata-server,
        gke.gcr.io/netd-amd64,
        gke.gcr.io/watcher-daemonset,
        gcr.io/google-containers/prometheus-to-sd,
        k8s.gcr.io/ip-masq-agent-amd64,
        k8s.gcr.io/kube-proxy,
        k8s.gcr.io/prometheus-to-sd,
        registry.k8s.io/ip-masq-agent-amd64,
        registry.k8s.io/kube-proxy,
        registry.k8s.io/prometheus-to-sd,
        quay.io/calico/node,
        sysdig/sysdig,
        sematext_images,
        k8s.gcr.io/dns/k8s-dns-node-cache,
        registry.k8s.io/dns/k8s-dns-node-cache,
        mcr.microsoft.com/oss/kubernetes/kube-proxy
      ]
    
    - macro: falco_privileged_containers
      condition: (openshift_image or
                  user_trusted_containers or
                  aws_eks_core_images or
                  container.image.repository in (trusted_images) or
                  container.image.repository in (falco_privileged_images) or
                  container.image.repository startswith istio/proxy_ or
                  container.image.repository startswith quay.io/sysdig/)
    
    # Add conditions to this macro (probably in a separate file,
    # overwriting this macro) to specify additional containers that are
    # allowed to run privileged
    #
    # In this file, it just takes one of the images in falco_privileged_images
    # and repeats it.
    - macro: user_privileged_containers
      condition: (never_true)
    
    # These container images are allowed to mount sensitive paths from the
    # host filesystem.
    - list: falco_sensitive_mount_images
      items: [
        falco_containers,
        docker.io/sysdig/sysdig, sysdig/sysdig,
        gcr.io/google_containers/hyperkube,
        gcr.io/google_containers/kube-proxy, docker.io/calico/node,
        docker.io/rook/toolbox, docker.io/cloudnativelabs/kube-router, docker.io/consul,
        docker.io/datadog/docker-dd-agent, docker.io/datadog/agent, docker.io/docker/ucp-agent, docker.io/gliderlabs/logspout,
        docker.io/netdata/netdata, docker.io/google/cadvisor, docker.io/prom/node-exporter,
        amazon/amazon-ecs-agent, prom/node-exporter, amazon/cloudwatch-agent
        ]
    
    - macro: falco_sensitive_mount_containers
      condition: (user_trusted_containers or
                  aws_eks_image_sensitive_mount or
                  container.image.repository in (trusted_images) or
                  container.image.repository in (falco_sensitive_mount_images) or
                  container.image.repository startswith quay.io/sysdig/ or
                  container.image.repository=k8scloudprovider/cinder-csi-plugin)
    
    # Add conditions to this macro (probably in a separate file,
    # overwriting this macro) to specify additional containers that are
    # allowed to perform sensitive mounts.
    #
    # In this file, it just takes one of the images in falco_sensitive_mount_images
    # and repeats it.
    - macro: user_sensitive_mount_containers
      condition: (never_true)
    
    - rule: Launch Privileged Container
      desc: Detect the initial process started in a privileged container. Exceptions are made for known trusted images.
      condition: >
        container_started and container
        and container.privileged=true
        and not falco_privileged_containers
        and not user_privileged_containers
        and not redhat_image
      output: Privileged container started (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid %container.info image=%container.image.repository:%container.image.tag)
      priority: INFO
      tags: [container, cis, mitre_privilege_escalation, mitre_lateral_movement, T1610]
    
    # These capabilities were used in the past to escape from containers
    - macro: excessively_capable_container
      condition: >
        (thread.cap_permitted contains CAP_SYS_ADMIN
        or thread.cap_permitted contains CAP_SYS_MODULE
        or thread.cap_permitted contains CAP_SYS_RAWIO
        or thread.cap_permitted contains CAP_SYS_PTRACE
        or thread.cap_permitted contains CAP_SYS_BOOT
        or thread.cap_permitted contains CAP_SYSLOG
        or thread.cap_permitted contains CAP_DAC_READ_SEARCH
        or thread.cap_permitted contains CAP_NET_ADMIN
        or thread.cap_permitted contains CAP_BPF)
    
    - rule: Launch Excessively Capable Container
      desc: Detect container started with a powerful set of capabilities. Exceptions are made for known trusted images.
      condition: >
        container_started and container
        and excessively_capable_container
        and not falco_privileged_containers
        and not user_privileged_containers
      output: Excessively capable container started (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid %container.info image=%container.image.repository:%container.image.tag cap_permitted=%thread.cap_permitted)
      priority: INFO
      tags: [container, cis, mitre_privilege_escalation, mitre_lateral_movement, T1610]
    
    
    # For now, only considering a full mount of /etc as
    # sensitive. Ideally, this would also consider all subdirectories
    # below /etc as well, but the globbing mechanism
    # doesn't allow exclusions of a full pattern, only single characters.
    - macro: sensitive_mount
      condition: (not container.mount.dest[/proc*] in ("<NA>","N/A") or
        not container.mount.dest[/var/run/docker.sock] in ("<NA>","N/A") or
        not container.mount.dest[/var/run/crio/crio.sock] in ("<NA>","N/A") or
        not container.mount.dest[/run/containerd/containerd.sock] in ("<NA>","N/A") or
        not container.mount.dest[/var/lib/kubelet] in ("<NA>","N/A") or
        not container.mount.dest[/var/lib/kubelet/pki] in ("<NA>","N/A") or
        not container.mount.dest[/] in ("<NA>","N/A") or
        not container.mount.dest[/home/admin] in ("<NA>","N/A") or
        not container.mount.dest[/etc] in ("<NA>","N/A") or
        not container.mount.dest[/etc/kubernetes] in ("<NA>","N/A") or
        not container.mount.dest[/etc/kubernetes/manifests] in ("<NA>","N/A") or
        not container.mount.dest[/root*] in ("<NA>","N/A"))
    
    # The steps libcontainer performs to set up the root program for a container are:
    # - clone + exec self to a program runc:[0:PARENT]
    # - clone a program runc:[1:CHILD] which sets up all the namespaces
    # - clone a second program runc:[2:INIT] + exec to the root program.
    #   The parent of runc:[2:INIT] is runc:0:PARENT]
    # As soon as 1:CHILD is created, 0:PARENT exits, so there's a race
    #   where at the time 2:INIT execs the root program, 0:PARENT might have
    #   already exited, or might still be around. So we handle both.
    # We also let runc:[1:CHILD] count as the parent process, which can occur
    # when we lose events and lose track of state.
    
    - macro: container_entrypoint
      condition: (not proc.pname exists or proc.pname in (runc:[0:PARENT], runc:[1:CHILD], runc, docker-runc, exe, docker-runc-cur))
    
    - rule: Launch Sensitive Mount Container
      desc: >
        Detect the initial process started by a container that has a mount from a sensitive host directory
        (i.e. /proc). Exceptions are made for known trusted images.
      condition: >
        container_started and container
        and sensitive_mount
        and not falco_sensitive_mount_containers
        and not user_sensitive_mount_containers
      output: Container with sensitive mount started (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid %container.info image=%container.image.repository:%container.image.tag mounts=%container.mounts)
      priority: INFO
      tags: [container, cis, mitre_lateral_movement, T1610]
    
    # In a local/user rules file, you could override this macro to
    # explicitly enumerate the container images that you want to run in
    # your environment. In this main falco rules file, there isn't any way
    # to know all the containers that can run, so any container is
    # allowed, by using a filter that is guaranteed to evaluate to true.
    # In the overridden macro, the condition would look something like
    # (container.image.repository = vendor/container-1 or
    # container.image.repository = vendor/container-2 or ...)
    - macro: allowed_containers
      condition: (container.id exists)
    
    - rule: Launch Disallowed Container
      desc: >
        Detect the initial process started by a container that is not in a list of allowed containers.
      condition: container_started and container and not allowed_containers
      output: Container started and not in allowed list (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid %container.info image=%container.image.repository:%container.image.tag)
      priority: WARNING
      tags: [container, mitre_lateral_movement, T1610]
    
    - macro: user_known_system_user_login
      condition: (never_true)
    
    # Anything run interactively by root
    # - condition: evt.type != switch and user.name = root and proc.name != sshd and interactive
    #  output: "Interactive root (%user.name %proc.name %evt.dir %evt.type %evt.args %fd.name)"
    #  priority: WARNING
    
    - rule: System user interactive
      desc: an attempt to run interactive commands by a system (i.e. non-login) user
      condition: spawned_process and system_users and interactive and not user_known_system_user_login
      output: "System user ran an interactive command (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid container_id=%container.id image=%container.image.repository exe_flags=%evt.arg.flags)"
      priority: INFO
      tags: [host, container, users, mitre_execution, T1059]
    
    # In some cases, a shell is expected to be run in a container. For example, configuration
    # management software may do this, which is expected.
    - macro: user_expected_terminal_shell_in_container_conditions
      condition: (never_true)
    
    - rule: Terminal shell in container
      desc: A shell was used as the entrypoint/exec point into a container with an attached terminal.
      condition: >
        spawned_process and container
        and shell_procs and proc.tty != 0
        and container_entrypoint
        and not user_expected_terminal_shell_in_container_conditions
      output: >
        A shell was spawned in a container with an attached terminal (user=%user.name user_loginuid=%user.loginuid %container.info
        shell=%proc.name parent=%proc.pname cmdline=%proc.cmdline pid=%proc.pid terminal=%proc.tty container_id=%container.id image=%container.image.repository exe_flags=%evt.arg.flags)
      priority: NOTICE
      tags: [container, shell, mitre_execution, T1059]
    
    # For some container types (mesos), there isn't a container image to
    # work with, and the container name is autogenerated, so there isn't
    # any stable aspect of the software to work with. In this case, we
    # fall back to allowing certain command lines.
    
    - list: known_shell_spawn_cmdlines
      items: [
        '"sh -c uname -p 2> /dev/null"',
        '"sh -c uname -s 2>&1"',
        '"sh -c uname -r 2>&1"',
        '"sh -c uname -v 2>&1"',
        '"sh -c uname -a 2>&1"',
        '"sh -c ruby -v 2>&1"',
        '"sh -c getconf CLK_TCK"',
        '"sh -c getconf PAGESIZE"',
        '"sh -c LC_ALL=C LANG=C /sbin/ldconfig -p 2>/dev/null"',
        '"sh -c LANG=C /sbin/ldconfig -p 2>/dev/null"',
        '"sh -c /sbin/ldconfig -p 2>/dev/null"',
        '"sh -c stty -a 2>/dev/null"',
        '"sh -c stty -a < /dev/tty"',
        '"sh -c stty -g < /dev/tty"',
        '"sh -c node index.js"',
        '"sh -c node index"',
        '"sh -c node ./src/start.js"',
        '"sh -c node app.js"',
        '"sh -c node -e \"require(''nan'')\""',
        '"sh -c node -e \"require(''nan'')\")"',
        '"sh -c node $NODE_DEBUG_OPTION index.js "',
        '"sh -c crontab -l 2"',
        '"sh -c lsb_release -a"',
        '"sh -c lsb_release -is 2>/dev/null"',
        '"sh -c whoami"',
        '"sh -c node_modules/.bin/bower-installer"',
        '"sh -c /bin/hostname -f 2> /dev/null"',
        '"sh -c locale -a"',
        '"sh -c  -t -i"',
        '"sh -c openssl version"',
        '"bash -c id -Gn kafadmin"',
        '"sh -c /bin/sh -c ''date +%%s''"',
        '"sh -c /usr/share/lighttpd/create-mime.conf.pl"'
        ]
    
    # This list allows for easy additions to the set of commands allowed
    # to run shells in containers without having to without having to copy
    # and override the entire run shell in container macro. Once
    # https://github.com/falcosecurity/falco/issues/255 is fixed this will be a
    # bit easier, as someone could append of any of the existing lists.
    - list: user_known_shell_spawn_binaries
      items: []
    
    # This macro allows for easy additions to the set of commands allowed
    # to run shells in containers without having to override the entire
    # rule. Its default value is an expression that always is false, which
    # becomes true when the "not ..." in the rule is applied.
    - macro: user_shell_container_exclusions
      condition: (never_true)
    
    - macro: login_doing_dns_lookup
      condition: (proc.name=login and fd.l4proto=udp and fd.sport=53)
    
    # sockfamily ip is to exclude certain processes (like 'groups') that communicate on unix-domain sockets
    # systemd can listen on ports to launch things like sshd on demand
    - rule: System procs network activity
      desc: any network activity performed by system binaries that are not expected to send or receive any network traffic
      condition: >
        (fd.sockfamily = ip and (system_procs or proc.name in (shell_binaries)))
        and (inbound_outbound)
        and not proc.name in (known_system_procs_network_activity_binaries)
        and not login_doing_dns_lookup
        and not user_expected_system_procs_network_activity_conditions
      output: >
        Known system binary sent/received network traffic
        (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid connection=%fd.name container_id=%container.id image=%container.image.repository)
      priority: NOTICE
      tags: [host, container, network, mitre_exfiltration, T1059, TA0011]
    
    # This list allows easily whitelisting system proc names that are
    # expected to communicate on the network.
    - list: known_system_procs_network_activity_binaries
      items: [systemd, hostid, id]
    
    # This macro allows specifying conditions under which a system binary
    # is allowed to communicate on the network. For instance, only specific
    # proc.cmdline values could be allowed to be more granular in what is
    # allowed.
    - macro: user_expected_system_procs_network_activity_conditions
      condition: (never_true)
    
    # When filled in, this should look something like:
    # (proc.env contains "HTTP_PROXY=http://my.http.proxy.com ")
    # The trailing space is intentional so avoid matching on prefixes of
    # the actual proxy.
    - macro: allowed_ssh_proxy_env
      condition: (always_true)
    
    - list: http_proxy_binaries
      items: [curl, wget]
    
    - macro: http_proxy_procs
      condition: (proc.name in (http_proxy_binaries))
    
    - rule: Program run with disallowed http proxy env
      desc: An attempt to run a program with a disallowed HTTP_PROXY environment variable
      condition: >
        spawned_process and
        http_proxy_procs and
        not allowed_ssh_proxy_env and
        proc.env icontains HTTP_PROXY
      enabled: false
      output: >
        Program run with disallowed HTTP_PROXY environment variable
        (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid env=%proc.env parent=%proc.pname container_id=%container.id image=%container.image.repository exe_flags=%evt.arg.flags)
      priority: NOTICE
      tags: [host, container, users, mitre_command_and_control, T1090, T1204]
    
    # In some environments, any attempt by a interpreted program (perl,
    # python, ruby, etc) to listen for incoming connections or perform
    # outgoing connections might be suspicious. These rules are not
    # enabled by default.
    
    - rule: Interpreted procs inbound network activity
      desc: Any inbound network activity performed by any interpreted program (perl, python, ruby, etc.)
      condition: >
        (inbound and interpreted_procs)
      enabled: false
      output: >
        Interpreted program received/listened for network traffic
        (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid connection=%fd.name container_id=%container.id image=%container.image.repository)
      priority: NOTICE
      tags: [host, container, network, mitre_exfiltration, TA0011]
    
    - rule: Interpreted procs outbound network activity
      desc: Any outbound network activity performed by any interpreted program (perl, python, ruby, etc.)
      condition: >
        (outbound and interpreted_procs)
      enabled: false
      output: >
        Interpreted program performed outgoing network connection
        (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid connection=%fd.name container_id=%container.id image=%container.image.repository)
      priority: NOTICE
      tags: [host, container, network, mitre_exfiltration, TA0011]
    
    - list: openvpn_udp_ports
      items: [1194, 1197, 1198, 8080, 9201]
    
    - list: l2tp_udp_ports
      items: [500, 1701, 4500, 10000]
    
    - list: statsd_ports
      items: [8125]
    
    - list: ntp_ports
      items: [123]
    
    # Some applications will connect a udp socket to an address only to
    # test connectivity. Assuming the udp connect works, they will follow
    # up with a tcp connect that actually sends/receives data.
    #
    # With that in mind, we listed a few commonly seen ports here to avoid
    # some false positives. In addition, we make the main rule opt-in, so
    # it's disabled by default.
    
    - list: test_connect_ports
      items: [0, 9, 80, 3306]
    
    - list: expected_udp_ports
      items: [53, openvpn_udp_ports, l2tp_udp_ports, statsd_ports, ntp_ports, test_connect_ports]
    
    - macro: expected_udp_traffic
      condition: fd.port in (expected_udp_ports)
    
    - rule: Unexpected UDP Traffic
      desc: UDP traffic not on port 53 (DNS) or other commonly used ports
      condition: (inbound_outbound) and fd.l4proto=udp and not expected_udp_traffic
      enabled: false
      output: >
        Unexpected UDP Traffic Seen
        (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid connection=%fd.name proto=%fd.l4proto evt=%evt.type %evt.args container_id=%container.id image=%container.image.repository)
      priority: NOTICE
      tags: [host, container, network, mitre_exfiltration, TA0011]
    
    # With the current restriction on system calls handled by falco
    # (e.g. excluding read/write/sendto/recvfrom/etc, this rule won't
    # trigger).
    # - rule: Ssh error in syslog
    #   desc: any ssh errors (failed logins, disconnects, ...) sent to syslog
    #   condition: syslog and ssh_error_message and evt.dir = <
    #   output: "sshd sent error message to syslog (error=%evt.buffer)"
    #   priority: WARNING
    
    - macro: somebody_becoming_themselves
      condition: ((user.name=nobody and evt.arg.uid=nobody) or
                  (user.name=www-data and evt.arg.uid=www-data) or
                  (user.name=_apt and evt.arg.uid=_apt) or
                  (user.name=postfix and evt.arg.uid=postfix) or
                  (user.name=pki-agent and evt.arg.uid=pki-agent) or
                  (user.name=pki-acme and evt.arg.uid=pki-acme) or
                  (user.name=nfsnobody and evt.arg.uid=nfsnobody) or
                  (user.name=postgres and evt.arg.uid=postgres))
    
    - macro: nrpe_becoming_nagios
      condition: (proc.name=nrpe and evt.arg.uid=nagios)
    
    # In containers, the user name might be for a uid that exists in the
    # container but not on the host. (See
    # https://github.com/draios/sysdig/issues/954). So in that case, allow
    # a setuid.
    - macro: known_user_in_container
      condition: (container and not user.name in ("<NA>","N/A",""))
    
    # Add conditions to this macro (probably in a separate file,
    # overwriting this macro) to allow for specific combinations of
    # programs changing users by calling setuid.
    #
    # In this file, it just takes one of the condition in the base macro
    # and repeats it.
    - macro: user_known_non_sudo_setuid_conditions
      condition: user.name=root
    
    # sshd, mail programs attempt to setuid to root even when running as non-root. Excluding here to avoid meaningless FPs
    - rule: Non sudo setuid
      desc: >
        an attempt to change users by calling setuid. sudo/su are excluded. users "root" and "nobody"
        suing to itself are also excluded, as setuid calls typically involve dropping privileges.
      condition: >
        evt.type=setuid and evt.dir=>
        and (known_user_in_container or not container)
        and not (user.name=root or user.uid=0)
        and not somebody_becoming_themselves
        and not proc.name in (known_setuid_binaries, userexec_binaries, mail_binaries, docker_binaries,
                              nomachine_binaries)
        and not proc.name startswith "runc:"
        and not java_running_sdjagent
        and not nrpe_becoming_nagios
        and not user_known_non_sudo_setuid_conditions
      output: >
        Unexpected setuid call by non-sudo, non-root program (user=%user.name user_loginuid=%user.loginuid cur_uid=%user.uid parent=%proc.pname
        command=%proc.cmdline pid=%proc.pid uid=%evt.arg.uid container_id=%container.id image=%container.image.repository)
      priority: NOTICE
      tags: [host, container, users, mitre_privilege_escalation, T1548.001]
    
    - macro: user_known_user_management_activities
      condition: (never_true)
    
    - macro: chage_list
      condition: (proc.name=chage and (proc.cmdline contains "-l" or proc.cmdline contains "--list"))
    
    - rule: User mgmt binaries
      desc: >
        activity by any programs that can manage users, passwords, or permissions. sudo and su are excluded.
        Activity in containers is also excluded--some containers create custom users on top
        of a base linux distribution at startup.
        Some innocuous command lines that don't actually change anything are excluded.
      condition: >
        spawned_process and proc.name in (user_mgmt_binaries) and
        not proc.name in (su, sudo, lastlog, nologin, unix_chkpwd) and not container and
        not proc.pname in (cron_binaries, systemd, systemd.postins, udev.postinst, run-parts) and
        not proc.cmdline startswith "passwd -S" and
        not proc.cmdline startswith "useradd -D" and
        not proc.cmdline startswith "systemd --version" and
        not run_by_qualys and
        not run_by_sumologic_securefiles and
        not run_by_yum and
        not run_by_ms_oms and
        not run_by_google_accounts_daemon and
        not chage_list and
        not user_known_user_management_activities
      output: >
        User management binary command run outside of container
        (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid parent=%proc.pname gparent=%proc.aname[2] ggparent=%proc.aname[3] gggparent=%proc.aname[4] exe_flags=%evt.arg.flags)
      priority: NOTICE
      tags: [host, container, users, software_mgmt, mitre_persistence, T1543, T1098]
    
    - list: allowed_dev_files
      items: [
        /dev/null, /dev/stdin, /dev/stdout, /dev/stderr,
        /dev/random, /dev/urandom, /dev/console, /dev/kmsg
        ]
    
    - macro: user_known_create_files_below_dev_activities
      condition: (never_true)
    
    # (we may need to add additional checks against false positives, see:
    # https://bugs.launchpad.net/ubuntu/+source/rkhunter/+bug/86153)
    - rule: Create files below dev
      desc: creating any files below /dev other than known programs that manage devices. Some rootkits hide files in /dev.
      condition: >
        fd.directory = /dev and
        (evt.type = creat or (evt.type in (open,openat,openat2) and evt.arg.flags contains O_CREAT))
        and not proc.name in (dev_creation_binaries)
        and not fd.name in (allowed_dev_files)
        and not fd.name startswith /dev/tty
        and not user_known_create_files_below_dev_activities
      output: "File created below /dev by untrusted program (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid file=%fd.name container_id=%container.id image=%container.image.repository)"
      priority: ERROR
      tags: [host, container, filesystem, mitre_persistence, T1543, T1083]
    
    
    # In a local/user rules file, you could override this macro to
    # explicitly enumerate the container images that you want to allow
    # access to EC2 metadata. In this main falco rules file, there isn't
    # any way to know all the containers that should have access, so any
    # container is allowed, by repeating the "container" macro. In the
    # overridden macro, the condition would look something like
    # (container.image.repository = vendor/container-1 or
    # container.image.repository = vendor/container-2 or ...)
    - macro: ec2_metadata_containers
      condition: container
    
    # On EC2 instances, 169.254.169.254 is a special IP used to fetch
    # metadata about the instance. It may be desirable to prevent access
    # to this IP from containers.
    - rule: Contact EC2 Instance Metadata Service From Container
      desc: Detect attempts to contact the EC2 Instance Metadata Service from a container
      condition: outbound and fd.sip="169.254.169.254" and container and not ec2_metadata_containers
      output: Outbound connection to EC2 instance metadata service (command=%proc.cmdline pid=%proc.pid connection=%fd.name %container.info image=%container.image.repository:%container.image.tag)
      priority: NOTICE
      enabled: false
      tags: [network, aws, container, mitre_discovery, T1565]
    
    
    # This rule is not enabled by default, since this rule is for cloud environment(GCP, AWS and Azure) only.
    # You can filter the container that you want to allow access to metadata by overwriting user_known_metadata_access macro.
    
    - macro: user_known_metadata_access
      condition: (k8s.ns.name = "kube-system")
    
    # On GCP, AWS and Azure, 169.254.169.254 is a special IP used to fetch
    # metadata about the instance. The metadata could be used to get credentials by attackers.
    - rule: Contact cloud metadata service from container
      desc: Detect attempts to contact the Cloud Instance Metadata Service from a container
      condition: outbound and fd.sip="169.254.169.254" and container and not user_known_metadata_access
      enabled: false
      output: Outbound connection to cloud instance metadata service (command=%proc.cmdline pid=%proc.pid connection=%fd.name %container.info image=%container.image.repository:%container.image.tag)
      priority: NOTICE
      tags: [network, container, mitre_discovery, T1565]
    
    # Containers from IBM Cloud
    - list: ibm_cloud_containers
      items:
        - icr.io/ext/sysdig/agent
        - registry.ng.bluemix.net/armada-master/metrics-server-amd64
        - registry.ng.bluemix.net/armada-master/olm
    
    # In a local/user rules file, list the namespace or container images that are
    # allowed to contact the K8s API Server from within a container. This
    # might cover cases where the K8s infrastructure itself is running
    # within a container.
    # TODO: Remove k8s.gcr.io reference after 01/Dec/2023
    - macro: k8s_containers
      condition: >
        (container.image.repository in (gcr.io/google_containers/hyperkube-amd64,
         gcr.io/google_containers/kube2sky,
         docker.io/sysdig/sysdig, sysdig/sysdig,
         fluent/fluentd-kubernetes-daemonset, prom/prometheus,
         falco_containers,
         falco_no_driver_containers,
         ibm_cloud_containers,
         velero/velero,
         quay.io/jetstack/cert-manager-cainjector, weaveworks/kured,
         quay.io/prometheus-operator/prometheus-operator, k8s.gcr.io/ingress-nginx/kube-webhook-certgen,
         registry.k8s.io/ingress-nginx/kube-webhook-certgen, quay.io/spotahome/redis-operator,
         registry.opensource.zalan.do/acid/postgres-operator, registry.opensource.zalan.do/acid/postgres-operator-ui,
         rabbitmqoperator/cluster-operator, quay.io/kubecost1/kubecost-cost-model,
         docker.io/bitnami/prometheus, docker.io/bitnami/kube-state-metrics, mcr.microsoft.com/oss/azure/aad-pod-identity/nmi)
         or (k8s.ns.name = "kube-system"))
    
    - macro: k8s_api_server
      condition: (fd.sip.name="kubernetes.default.svc.cluster.local")
    
    - macro: user_known_contact_k8s_api_server_activities
      condition: (never_true)
    
    - rule: Contact K8S API Server From Container
      desc: Detect attempts to contact the K8S API Server from a container
      condition: >
        evt.type=connect and evt.dir=< and
        (fd.typechar=4 or fd.typechar=6) and
        container and
        not k8s_containers and
        k8s_api_server and
        not user_known_contact_k8s_api_server_activities
      output: Unexpected connection to K8s API Server from container (command=%proc.cmdline pid=%proc.pid %container.info image=%container.image.repository:%container.image.tag connection=%fd.name)
      priority: NOTICE
      tags: [network, k8s, container, mitre_discovery, T1565]
    
    # In a local/user rules file, list the container images that are
    # allowed to contact NodePort services from within a container. This
    # might cover cases where the K8s infrastructure itself is running
    # within a container.
    #
    # By default, all containers are allowed to contact NodePort services.
    - macro: nodeport_containers
      condition: container
    
    - rule: Unexpected K8s NodePort Connection
      desc: Detect attempts to use K8s NodePorts from a container
      condition: (inbound_outbound) and fd.sport >= 30000 and fd.sport <= 32767 and container and not nodeport_containers
      output: Unexpected K8s NodePort Connection (command=%proc.cmdline pid=%proc.pid connection=%fd.name container_id=%container.id image=%container.image.repository)
      priority: NOTICE
      tags: [network, k8s, container, mitre_persistence, T1205.001]
    
    - list: network_tool_binaries
      items: [nc, ncat, netcat, nmap, dig, tcpdump, tshark, ngrep, telnet, mitmproxy, socat, zmap]
    
    - macro: network_tool_procs
      condition: (proc.name in (network_tool_binaries))
    
    # In a local/user rules file, create a condition that matches legitimate uses
    # of a package management process inside a container.
    #
    # For example:
    # - macro: user_known_package_manager_in_container
    #   condition: proc.cmdline="dpkg -l"
    - macro: user_known_package_manager_in_container
      condition: (never_true)
    
    # Container is supposed to be immutable. Package management should be done in building the image.
    # TODO: Remove k8s.gcr.io reference after 01/Dec/2023
    - macro: pkg_mgmt_in_kube_proxy
      condition: >
        proc.cmdline startswith "update-alternat"
        and (container.image.repository = "registry.k8s.io/kube-proxy"
        or container.image.repository = "k8s.gcr.io/kube-proxy")
    
    - rule: Launch Package Management Process in Container
      desc: Package management process ran inside container
      condition: >
        spawned_process
        and container
        and user.name != "_apt"
        and package_mgmt_procs
        and not package_mgmt_ancestor_procs
        and not user_known_package_manager_in_container
        and not pkg_mgmt_in_kube_proxy
      output: >
        Package management process launched in container (user=%user.name user_loginuid=%user.loginuid
        command=%proc.cmdline pid=%proc.pid container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority: ERROR
      tags: [container, process, software_mgmt, mitre_persistence, T1505]
    
    - rule: Netcat Remote Code Execution in Container
      desc: Netcat Program runs inside container that allows remote code execution
      condition: >
        spawned_process and container and
        ((proc.name = "nc" and (proc.args contains "-e" or proc.args contains "-c")) or
         (proc.name = "ncat" and (proc.args contains "--sh-exec" or proc.args contains "--exec" or proc.args contains "-e "
                                  or proc.args contains "-c " or proc.args contains "--lua-exec"))
        )
      output: >
        Netcat runs inside container that allows remote code execution (user=%user.name user_loginuid=%user.loginuid
        command=%proc.cmdline pid=%proc.pid container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority: WARNING
      tags: [container, network, process, mitre_execution, T1059]
    
    - macro: user_known_network_tool_activities
      condition: (never_true)
    
    - rule: Launch Suspicious Network Tool in Container
      desc: Detect network tools launched inside container
      condition: >
        spawned_process and container and network_tool_procs and not user_known_network_tool_activities
      output: >
        Network tool launched in container (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid parent_process=%proc.pname
        container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority: NOTICE
      tags: [container, network, process, mitre_discovery, mitre_exfiltration, T1595, T1046]
    
    # This rule is not enabled by default, as there are legitimate use
    # cases for these tools on hosts. If you want to enable it, modify the
    # following macro.
    - macro: consider_network_tools_on_host
      condition: (never_true)
    
    - rule: Launch Suspicious Network Tool on Host
      desc: Detect network tools launched on the host
      condition: >
        spawned_process and
        not container and
        consider_network_tools_on_host and
        network_tool_procs and
        not user_known_network_tool_activities
      output: >
        Network tool launched on host (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid parent_process=%proc.pname exe_flags=%evt.arg.flags)
      priority: NOTICE
      tags: [host, network, process, mitre_discovery, mitre_exfiltration, T1595, T1046]
    
    - list: grep_binaries
      items: [grep, egrep, fgrep]
    
    - macro: grep_commands
      condition: (proc.name in (grep_binaries))
    
    # a less restrictive search for things that might be passwords/ssh/user etc.
    - macro: grep_more
      condition: (never_true)
    
    - macro: private_key_or_password
      condition: >
        (proc.args icontains "BEGIN PRIVATE" or
         proc.args icontains "BEGIN OPENSSH PRIVATE" or
         proc.args icontains "BEGIN RSA PRIVATE" or
         proc.args icontains "BEGIN DSA PRIVATE" or
         proc.args icontains "BEGIN EC PRIVATE" or
         (grep_more and
          (proc.args icontains " pass " or
           proc.args icontains " ssh " or
           proc.args icontains " user "))
        )
    
    - rule: Search Private Keys or Passwords
      desc: >
        Detect grep private keys or passwords activity.
      condition: >
        (spawned_process and
         ((grep_commands and private_key_or_password) or
          (proc.name = "find" and
            (proc.args contains "id_rsa" or 
             proc.args contains "id_dsa" or 
             proc.args contains "id_ed25519" or 
             proc.args contains "id_ecdsa"
            )
          ))
        )
      output: >
        Grep private keys or passwords activities found
        (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid container_id=%container.id container_name=%container.name
        image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority:
        WARNING
      tags: [host, container, process, filesystem, mitre_credential_access, T1552.001]
    
    - list: log_directories
      items: [/var/log, /dev/log]
    
    - list: log_files
      items: [syslog, auth.log, secure, kern.log, cron, user.log, dpkg.log, last.log, yum.log, access_log, mysql.log, mysqld.log]
    
    - macro: access_log_files
      condition: (fd.directory in (log_directories) or fd.filename in (log_files))
    
    # a placeholder for whitelist log files that could be cleared. Recommend the macro as (fd.name startswith "/var/log/app1*")
    - macro: allowed_clear_log_files
      condition: (never_true)
    
    - macro: trusted_logging_images
      condition: (container.image.repository endswith "splunk/fluentd-hec" or
                  container.image.repository endswith "fluent/fluentd-kubernetes-daemonset" or
                  container.image.repository endswith "openshift3/ose-logging-fluentd" or
                  container.image.repository endswith "containernetworking/azure-npm")
    
    - rule: Clear Log Activities
      desc: Detect clearing of critical log files
      condition: >
        open_write and
        access_log_files and
        evt.arg.flags contains "O_TRUNC" and
        not trusted_logging_images and
        not allowed_clear_log_files
      output: >
        Log files were tampered (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid file=%fd.name container_id=%container.id image=%container.image.repository)
      priority:
        WARNING
      tags: [host, container, filesystem, mitre_defense_evasion, T1070]
    
    - list: data_remove_commands
      items: [shred, mkfs, mke2fs]
    
    - macro: clear_data_procs
      condition: (proc.name in (data_remove_commands))
    
    - macro: user_known_remove_data_activities
      condition: (never_true)
    
    - rule: Remove Bulk Data from Disk
      desc: Detect process running to clear bulk data from disk
      condition: spawned_process and clear_data_procs and not user_known_remove_data_activities
      output: >
        Bulk data has been removed from disk (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid file=%fd.name container_id=%container.id image=%container.image.repository exe_flags=%evt.arg.flags)
      priority:
        WARNING
      tags: [host, container, process, filesystem, mitre_persistence, T1485]
    
    # here `+"`"+`ash_history`+"`"+` will match both `+"`"+`bash_history`+"`"+` and `+"`"+`ash_history`+"`"+`
    - macro: modify_shell_history
      condition: >
        (modify and (
          evt.arg.name endswith "ash_history" or
          evt.arg.name endswith "zsh_history" or
          evt.arg.name contains "fish_read_history" or
          evt.arg.name endswith "fish_history" or
          evt.arg.oldpath endswith "ash_history" or
          evt.arg.oldpath endswith "zsh_history" or
          evt.arg.oldpath contains "fish_read_history" or
          evt.arg.oldpath endswith "fish_history" or
          evt.arg.path endswith "ash_history" or
          evt.arg.path endswith "zsh_history" or
          evt.arg.path contains "fish_read_history" or
          evt.arg.path endswith "fish_history"))
    
    # here `+"`"+`ash_history`+"`"+` will match both `+"`"+`bash_history`+"`"+` and `+"`"+`ash_history`+"`"+`
    - macro: truncate_shell_history
      condition: >
        (open_write and (
          fd.name endswith "ash_history" or
          fd.name endswith "zsh_history" or
          fd.name contains "fish_read_history" or
          fd.name endswith "fish_history") and evt.arg.flags contains "O_TRUNC")
    
    - macro: var_lib_docker_filepath
      condition: (evt.arg.name startswith /var/lib/docker or fd.name startswith /var/lib/docker)
    
    - rule: Delete or rename shell history
      desc: Detect shell history deletion
      condition: >
        (modify_shell_history or truncate_shell_history) and
           not var_lib_docker_filepath and
           not proc.name in (docker_binaries)
      output: >
        Shell history had been deleted or renamed (user=%user.name user_loginuid=%user.loginuid type=%evt.type command=%proc.cmdline pid=%proc.pid fd.name=%fd.name name=%evt.arg.name path=%evt.arg.path oldpath=%evt.arg.oldpath %container.info)
      priority:
        WARNING
      tags: [host, container, process, filesystem, mitre_defense_evasion, T1070]
    
    # This rule is deprecated and will/should never be triggered. Keep it here for backport compatibility.
    # Rule Delete or rename shell history is the preferred rule to use now.
    - rule: Delete Bash History
      desc: Detect bash history deletion
      condition: >
        ((spawned_process and proc.name in (shred, rm, mv) and proc.args contains "bash_history") or
         (open_write and fd.name contains "bash_history" and evt.arg.flags contains "O_TRUNC"))
      output: >
        Shell history had been deleted or renamed (user=%user.name user_loginuid=%user.loginuid type=%evt.type command=%proc.cmdline pid=%proc.pid fd.name=%fd.name name=%evt.arg.name path=%evt.arg.path oldpath=%evt.arg.oldpath exe_flags=%evt.arg.flags %container.info)
      priority:
        WARNING
      tags: [host, container, process, filesystem, mitre_defense_evasion, T1070]
    
    - list: user_known_chmod_applications
      items: [hyperkube, kubelet, k3s-agent]
    
    # This macro should be overridden in user rules as needed. This is useful if a given application
    # should not be ignored altogether with the user_known_chmod_applications list, but only in
    # specific conditions.
    - macro: user_known_set_setuid_or_setgid_bit_conditions
      condition: (never_true)
    
    - rule: Set Setuid or Setgid bit
      desc: >
        When the setuid or setgid bits are set for an application,
        this means that the application will run with the privileges of the owning user or group respectively.
        Detect setuid or setgid bits set via chmod
      condition: >
        chmod and (evt.arg.mode contains "S_ISUID" or evt.arg.mode contains "S_ISGID")
        and not proc.name in (user_known_chmod_applications)
        and not exe_running_docker_save
        and not user_known_set_setuid_or_setgid_bit_conditions
      enabled: false
      output: >
        Setuid or setgid bit is set via chmod (fd=%evt.arg.fd filename=%evt.arg.filename mode=%evt.arg.mode user=%user.name user_loginuid=%user.loginuid process=%proc.name
        command=%proc.cmdline pid=%proc.pid container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag)
      priority:
        NOTICE
      tags: [host, container, process, users, mitre_persistence, T1548.001]
    
    - list: exclude_hidden_directories
      items: [/root/.cassandra]
    
    # The rule is disabled by default.
    - macro: user_known_create_hidden_file_activities
      condition: (never_true)
    
    - rule: Create Hidden Files or Directories
      desc: Detect hidden files or directories created
      condition: >
        ((modify and evt.arg.newpath contains "/.") or
         (mkdir and evt.arg.path contains "/.") or
         (open_write and evt.arg.flags contains "O_CREAT" and fd.name contains "/." and not fd.name pmatch (exclude_hidden_directories))) and
        not user_known_create_hidden_file_activities
        and not exe_running_docker_save
      enabled: false
      output: >
        Hidden file or directory created (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid
        file=%fd.name newpath=%evt.arg.newpath container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag)
      priority:
        NOTICE
      tags: [host, container, filesystem, mitre_persistence, T1564.001]
    
    - list: remote_file_copy_binaries
      items: [rsync, scp, sftp, dcp]
    
    - macro: remote_file_copy_procs
      condition: (proc.name in (remote_file_copy_binaries))
    
    # Users should overwrite this macro to specify conditions under which a
    # Custom condition for use of remote file copy tool in container
    - macro: user_known_remote_file_copy_activities
      condition: (never_true)
    
    - rule: Launch Remote File Copy Tools in Container
      desc: Detect remote file copy tools launched in container
      condition: >
        spawned_process
        and container
        and remote_file_copy_procs
        and not user_known_remote_file_copy_activities
      output: >
        Remote file copy tool launched in container (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid parent_process=%proc.pname
        container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority: NOTICE
      tags: [container, network, process, mitre_lateral_movement, mitre_exfiltration, T1020, T1210]
    
    - rule: Create Symlink Over Sensitive Files
      desc: Detect symlink created over sensitive files
      condition: >
        create_symlink and
        (evt.arg.target in (sensitive_file_names) or evt.arg.target in (sensitive_directory_names))
      output: >
        Symlinks created over sensitive files (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid target=%evt.arg.target linkpath=%evt.arg.linkpath parent_process=%proc.pname)
      priority: WARNING
      tags: [host, container, filesystem, mitre_exfiltration, mitre_credential_access, T1020, T1083, T1212, T1552, T1555]
    
    - rule: Create Hardlink Over Sensitive Files
      desc: Detect hardlink created over sensitive files
      condition: >
        create_hardlink and
        (evt.arg.oldpath in (sensitive_file_names))
      output: >
        Hardlinks created over sensitive files (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid target=%evt.arg.oldpath linkpath=%evt.arg.newpath parent_process=%proc.pname)
      priority: WARNING
      tags: [host, container, filesystem, mitre_exfiltration, mitre_credential_access, T1020, T1083, T1212, T1552, T1555]
    
    - list: miner_ports
      items: [
            25, 3333, 3334, 3335, 3336, 3357, 4444,
            5555, 5556, 5588, 5730, 6099, 6666, 7777,
            7778, 8000, 8001, 8008, 8080, 8118, 8333,
            8888, 8899, 9332, 9999, 14433, 14444,
            45560, 45700
        ]
    
    - list: miner_domains
      items: [
          "asia1.ethpool.org","ca.minexmr.com",
          "cn.stratum.slushpool.com","de.minexmr.com",
          "eth-ar.dwarfpool.com","eth-asia.dwarfpool.com",
          "eth-asia1.nanopool.org","eth-au.dwarfpool.com",
          "eth-au1.nanopool.org","eth-br.dwarfpool.com",
          "eth-cn.dwarfpool.com","eth-cn2.dwarfpool.com",
          "eth-eu.dwarfpool.com","eth-eu1.nanopool.org",
          "eth-eu2.nanopool.org","eth-hk.dwarfpool.com",
          "eth-jp1.nanopool.org","eth-ru.dwarfpool.com",
          "eth-ru2.dwarfpool.com","eth-sg.dwarfpool.com",
          "eth-us-east1.nanopool.org","eth-us-west1.nanopool.org",
          "eth-us.dwarfpool.com","eth-us2.dwarfpool.com",
          "eu.stratum.slushpool.com","eu1.ethermine.org",
          "eu1.ethpool.org","fr.minexmr.com",
          "mine.moneropool.com","mine.xmrpool.net",
          "pool.minexmr.com","pool.monero.hashvault.pro",
          "pool.supportxmr.com","sg.minexmr.com",
          "sg.stratum.slushpool.com","stratum-eth.antpool.com",
          "stratum-ltc.antpool.com","stratum-zec.antpool.com",
          "stratum.antpool.com","us-east.stratum.slushpool.com",
          "us1.ethermine.org","us1.ethpool.org",
          "us2.ethermine.org","us2.ethpool.org",
          "xmr-asia1.nanopool.org","xmr-au1.nanopool.org",
          "xmr-eu1.nanopool.org","xmr-eu2.nanopool.org",
          "xmr-jp1.nanopool.org","xmr-us-east1.nanopool.org",
          "xmr-us-west1.nanopool.org","xmr.crypto-pool.fr",
          "xmr.pool.minergate.com", "rx.unmineable.com",
          "ss.antpool.com","dash.antpool.com",
          "eth.antpool.com","zec.antpool.com",
          "xmc.antpool.com","btm.antpool.com",
          "stratum-dash.antpool.com","stratum-xmc.antpool.com",
          "stratum-btm.antpool.com"
          ]
    
    - list: https_miner_domains
      items: [
        "ca.minexmr.com",
        "cn.stratum.slushpool.com",
        "de.minexmr.com",
        "fr.minexmr.com",
        "mine.moneropool.com",
        "mine.xmrpool.net",
        "pool.minexmr.com",
        "sg.minexmr.com",
        "stratum-eth.antpool.com",
        "stratum-ltc.antpool.com",
        "stratum-zec.antpool.com",
        "stratum.antpool.com",
        "xmr.crypto-pool.fr",
        "ss.antpool.com",
        "stratum-dash.antpool.com",
        "stratum-xmc.antpool.com",
        "stratum-btm.antpool.com",
        "btm.antpool.com"
      ]
    
    - list: http_miner_domains
      items: [
        "ca.minexmr.com",
        "de.minexmr.com",
        "fr.minexmr.com",
        "mine.moneropool.com",
        "mine.xmrpool.net",
        "pool.minexmr.com",
        "sg.minexmr.com",
        "xmr.crypto-pool.fr"
      ]
    
    # Add rule based on crypto mining IOCs
    - macro: minerpool_https
      condition: (fd.sport="443" and fd.sip.name in (https_miner_domains))
    
    - macro: minerpool_http
      condition: (fd.sport="80" and fd.sip.name in (http_miner_domains))
    
    - macro: minerpool_other
      condition: (fd.sport in (miner_ports) and fd.sip.name in (miner_domains))
    
    - macro: net_miner_pool
      condition: (evt.type in (sendto, sendmsg, connect) and evt.dir=< and (fd.net != "127.0.0.0/8" and not fd.snet in (rfc_1918_addresses)) and ((minerpool_http) or (minerpool_https) or (minerpool_other)))
    
    - macro: trusted_images_query_miner_domain_dns
      condition: (container.image.repository in (falco_containers))
    
    # The rule is disabled by default.
    # Note: falco will send DNS request to resolve miner pool domain which may trigger alerts in your environment.
    - rule: Detect outbound connections to common miner pool ports
      desc: Miners typically connect to miner pools on common ports.
      condition: net_miner_pool and not trusted_images_query_miner_domain_dns
      enabled: false
      output: Outbound connection to IP/Port flagged by https://cryptoioc.ch (command=%proc.cmdline pid=%proc.pid port=%fd.rport ip=%fd.rip container=%container.info image=%container.image.repository)
      priority: CRITICAL
      tags: [host, container, network, mitre_execution, T1496]
    
    - rule: Detect crypto miners using the Stratum protocol
      desc: Miners typically specify the mining pool to connect to with a URI that begins with 'stratum+tcp'
      condition: spawned_process and (proc.cmdline contains "stratum+tcp" or proc.cmdline contains "stratum2+tcp" or proc.cmdline contains "stratum+ssl" or proc.cmdline contains "stratum2+ssl")
      output: Possible miner running (command=%proc.cmdline pid=%proc.pid container=%container.info image=%container.image.repository exe_flags=%evt.arg.flags)
      priority: CRITICAL
      tags: [host, container, process, mitre_execution, T1496]
    
    - list: k8s_client_binaries
      items: [docker, kubectl, crictl]
    
    # TODO: Remove k8s.gcr.io reference after 01/Dec/2023
    - list: user_known_k8s_ns_kube_system_images
      items: [
        k8s.gcr.io/fluentd-gcp-scaler,
        k8s.gcr.io/node-problem-detector/node-problem-detector,
        registry.k8s.io/fluentd-gcp-scaler,
        registry.k8s.io/node-problem-detector/node-problem-detector
      ]
    
    - list: user_known_k8s_images
      items: [
        mcr.microsoft.com/aks/hcp/hcp-tunnel-front
      ]
    
    # Whitelist for known docker client binaries run inside container
    # - k8s.gcr.io/fluentd-gcp-scaler / registry.k8s.io/fluentd-gcp-scaler in GCP/GKE
    # TODO: Remove k8s.gcr.io reference after 01/Dec/2023
    - macro: user_known_k8s_client_container
      condition: >
        (k8s.ns.name="kube-system" and container.image.repository in (user_known_k8s_ns_kube_system_images)) or container.image.repository in (user_known_k8s_images)
    
    - macro: user_known_k8s_client_container_parens
      condition: (user_known_k8s_client_container)
    
    - rule: The docker client is executed in a container
      desc: Detect a k8s client tool executed inside a container
      condition: spawned_process and container and not user_known_k8s_client_container_parens and proc.name in (k8s_client_binaries)
      output: "Docker or kubernetes client executed in container (user=%user.name user_loginuid=%user.loginuid %container.info parent=%proc.pname cmdline=%proc.cmdline pid=%proc.pid image=%container.image.repository:%container.image.tag)"
      priority: WARNING
      tags: [container, mitre_execution, T1610]
    
    - list: user_known_packet_socket_binaries
      items: []
    
    - rule: Packet socket created in container
      desc: Detect new packet socket at the device driver (OSI Layer 2) level in a container. Packet socket could be used for ARP Spoofing and privilege escalation(CVE-2020-14386) by attacker.
      condition: evt.type=socket and evt.arg[0] contains AF_PACKET and container and not proc.name in (user_known_packet_socket_binaries)
      output: Packet socket was created in a container (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid socket_info=%evt.args container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag)
      priority: NOTICE
      tags: [container, network, mitre_discovery, T1046]
    
    # Namespaces where the rule is enforce
    - list: namespace_scope_network_only_subnet
      items: []
    
    - macro: network_local_subnet
      condition: >
        fd.rnet in (rfc_1918_addresses) or
        fd.ip = "0.0.0.0" or
        fd.net = "127.0.0.0/8"
    
    # # The rule is disabled by default.
    # # How to test:
    # # Add 'default' to namespace_scope_network_only_subnet
    # # Run:
    # kubectl run --generator=run-pod/v1 -n default -i --tty busybox --image=busybox --rm -- wget google.com -O /var/google.html
    # # Check logs running
    
    - rule: Network Connection outside Local Subnet
      desc: Detect traffic to image outside local subnet.
      condition: >
        inbound_outbound and
        container and
        not network_local_subnet and
        k8s.ns.name in (namespace_scope_network_only_subnet)
      enabled: false
      output: >
        Network connection outside local subnet
        (command=%proc.cmdline pid=%proc.pid connection=%fd.name user=%user.name user_loginuid=%user.loginuid container_id=%container.id
         image=%container.image.repository namespace=%k8s.ns.name
         fd.rip.name=%fd.rip.name fd.lip.name=%fd.lip.name fd.cip.name=%fd.cip.name fd.sip.name=%fd.sip.name)
      priority: WARNING
      tags: [container, network, mitre_discovery, T1046]
    
    - list: allowed_image
      items: [] # add image to monitor, i.e.: bitnami/nginx
    
    - list: authorized_server_binary
      items: []  # add binary to allow, i.e.: nginx
      
    - list: authorized_server_port
      items: [] # add port to allow, i.e.: 80
    
    # # How to test:
    # kubectl run --image=nginx nginx-app --port=80 --env="DOMAIN=cluster"
    # kubectl expose deployment nginx-app --port=80 --name=nginx-http --type=LoadBalancer
    # # On minikube:
    # minikube service nginx-http
    # # On general K8s:
    # kubectl get services
    # kubectl cluster-info
    # # Visit the Nginx service and port, should not fire.
    # # Change rule to different port, then different process name, and test again that it fires.
    
    - rule: Outbound or Inbound Traffic not to Authorized Server Process and Port
      desc: Detect traffic that is not to authorized server process and port.
      condition: >
        inbound_outbound and
        container and
        container.image.repository in (allowed_image) and
        not proc.name in (authorized_server_binary) and
        not fd.sport in (authorized_server_port)
      enabled: false
      output: >
        Network connection outside authorized port and binary
        (command=%proc.cmdline pid=%proc.pid connection=%fd.name user=%user.name user_loginuid=%user.loginuid container_id=%container.id
        image=%container.image.repository)
      priority: WARNING
      tags: [container, network, mitre_discovery, TA0011]
    
    - macro: user_known_stand_streams_redirect_activities
      condition: (never_true)
    
    - macro: dup
      condition: evt.type in (dup, dup2, dup3)
    
    - rule: Redirect STDOUT/STDIN to Network Connection in Container
      desc: Detect redirecting stdout/stdin to network connection in container (potential reverse shell).
      condition: dup and container and evt.rawres in (0, 1, 2) and fd.type in ("ipv4", "ipv6") and not user_known_stand_streams_redirect_activities
      output: >
        Redirect stdout/stdin to network connection (user=%user.name user_loginuid=%user.loginuid %container.info process=%proc.name parent=%proc.pname cmdline=%proc.cmdline pid=%proc.pid terminal=%proc.tty container_id=%container.id image=%container.image.repository fd.name=%fd.name fd.num=%fd.num fd.type=%fd.type fd.sip=%fd.sip)
      priority: NOTICE
      tags: [container, network, process, mitre_discovery, mitre_execution, T1059]
    
    # The two Container Drift rules below will fire when a new executable is created in a container.
    # There are two ways to create executables - file is created with execution permissions or permissions change of existing file.
    # We will use a new filter, is_open_exec, to find all files creations with execution permission, and will trace all chmods in a container.
    # The use case we are targeting here is an attempt to execute code that was not shipped as part of a container (drift) -
    # an activity that might be malicious or non-compliant.
    # Two things to pay attention to:
    #   1) In most cases, 'docker cp' will not be identified, but the assumption is that if an attacker gained access to the container runtime daemon, they are already privileged
    #   2) Drift rules will be noisy in environments in which containers are built (e.g. docker build)
    # These two rules are not enabled by default.
    
    - macro: user_known_container_drift_activities
      condition: (always_true)
    
    - rule: Container Drift Detected (chmod)
      desc: New executable created in a container due to chmod
      condition: >
        chmod and
        container and
        not runc_writing_exec_fifo and
        not runc_writing_var_lib_docker and
        not user_known_container_drift_activities and
        evt.rawres>=0 and
        ((evt.arg.mode contains "S_IXUSR") or
        (evt.arg.mode contains "S_IXGRP") or
        (evt.arg.mode contains "S_IXOTH"))
      enabled: false
      output: Drift detected (chmod), new executable created in a container (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid filename=%evt.arg.filename name=%evt.arg.name mode=%evt.arg.mode event=%evt.type)
      priority: ERROR
      tags: [container, process, filesystem, mitre_execution, T1059]
    
    # ****************************************************************************
    # * "Container Drift Detected (open+create)" requires FALCO_ENGINE_VERSION 6 *
    # ****************************************************************************
    - rule: Container Drift Detected (open+create)
      desc: New executable created in a container due to open+create
      condition: >
        evt.type in (open,openat,openat2,creat) and
        evt.is_open_exec=true and
        container and
        not runc_writing_exec_fifo and
        not runc_writing_var_lib_docker and
        not user_known_container_drift_activities and
        evt.rawres>=0
      enabled: false
      output: Drift detected (open+create), new executable created in a container (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid filename=%evt.arg.filename name=%evt.arg.name mode=%evt.arg.mode event=%evt.type)
      priority: ERROR
      tags: [container, process, filesystem, mitre_execution, T1059]
    
    - list: c2_server_ip_list
      items: []
    
    - list: c2_server_fqdn_list
      items: []
    
    # For example, you can fetch a list of IP addresses and FQDN on this website:
    # https://feodotracker.abuse.ch/downloads/ipblocklist_recommended.json.
    # Use Falco HELM chart to update (append) the c2 server lists with your values.
    # See an example below.
    #
    #  `+"`"+``+"`"+``+"`"+`yaml
    #  # values.yaml Falco HELM chart file
    #  [...]
    #  customRules:
    #    c2-servers-list.yaml: |-
    #      - list: c2_server_ip_list
    #        append: true
    #        items: 
    #        - "'51.178.161.32'"
    #        - "'46.101.90.205'"
    #        
    #      - list: c2_server_fqdn_list
    #        append: true
    #        items: 
    #        - "srv-web.ffconsulting.com"
    #        - "57.ip-142-44-247.net"
    #  `+"`"+``+"`"+``+"`"+`
    
    - rule: Outbound Connection to C2 Servers
      desc: Detect outbound connection to command & control servers thanks to a list of IP addresses & a list of FQDN.
      condition: >
        outbound and 
        ((fd.sip in (c2_server_ip_list)) or
         (fd.sip.name in (c2_server_fqdn_list)))
      output: Outbound connection to C2 server (c2_domain=%fd.sip.name c2_addr=%fd.sip command=%proc.cmdline connection=%fd.name user=%user.name user_loginuid=%user.loginuid container_id=%container.id image=%container.image.repository)
      priority: WARNING
      enabled: false
      tags: [host, container, network, mitre_command_and_control, TA0011]
    
    - list: allowed_container_images_loading_kernel_module
      items: []
    
    # init_module and finit_module syscalls are available since Falco 0.35.0
    # rule coverage now extends to modprobe usage via init_module logging
    # and previous alerting on spawned_process and insmod is now covered
    # by finit_module syscall
    - rule: Linux Kernel Module Injection Detected
      desc: Detect kernel module was injected (from container).
      condition: kernel_module_load and container
        and not container.image.repository in (allowed_container_images_loading_kernel_module)
        and thread.cap_effective icontains sys_module
      output: Linux Kernel Module injection from container detected (user=%user.name uid=%user.uid user_loginuid=%user.loginuid process_name=%proc.name parent_process_name=%proc.pname parent_exepath=%proc.pexepath %proc.aname[2] %proc.aexepath[2] module=%proc.args %container.info image=%container.image.repository:%container.image.tag res=%evt.res syscall=%evt.type)
      priority: WARNING
      tags: [host, container, process, mitre_execution, mitre_persistence, TA0002]
    
    - list: run_as_root_image_list
      items: []
    
    - macro: user_known_run_as_root_container
      condition: (container.image.repository in (run_as_root_image_list))
    
    # The rule is disabled by default and should be enabled when non-root container policy has been applied.
    # Note the rule will not work as expected when usernamespace is applied, e.g. userns-remap is enabled.
    - rule: Container Run as Root User
      desc: Detected container running as root user
      condition: spawned_process and container and proc.vpid=1 and user.uid=0 and not user_known_run_as_root_container
      enabled: false
      output: Container launched with root user privilege (uid=%user.uid container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority: INFO
      tags: [container, process, users, mitre_execution, T1610]
    
    # This rule helps detect CVE-2021-3156:
    # A privilege escalation to root through heap-based buffer overflow
    - rule: Sudo Potential Privilege Escalation
      desc: Privilege escalation vulnerability affecting sudo (<= 1.9.5p2). Executing sudo using sudoedit -s or sudoedit -i command with command-line argument that ends with a single backslash character from an unprivileged user it's possible to elevate the user privileges to root.
      condition: spawned_process and user.uid != 0 and (proc.name=sudoedit or proc.name = sudo) and (proc.args contains -s or proc.args contains -i or proc.args contains --login) and (proc.args contains "\ " or proc.args endswith \)
      output: "Detect Sudo Privilege Escalation Exploit (CVE-2021-3156) (user=%user.name parent=%proc.pname cmdline=%proc.cmdline pid=%proc.pid exe_flags=%evt.arg.flags %container.info)"
      priority: CRITICAL
      tags: [host, container, filesystem, users, mitre_privilege_escalation, T1548.003]
    
    - rule: Debugfs Launched in Privileged Container
      desc: Detect file system debugger debugfs launched inside a privileged container which might lead to container escape.
      condition: >
        spawned_process and container
        and container.privileged=true
        and proc.name=debugfs
      output: Debugfs launched started in a privileged container (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid %container.info image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority: WARNING
      tags: [container, cis, process, mitre_execution, mitre_lateral_movement, T1611]
    
    - macro: mount_info
      condition: (proc.args="" or proc.args intersects ("-V", "-l", "-h"))
    
    - macro: known_gke_mount_in_privileged_containers
      condition:
        (k8s.ns.name = kube-system
        and container.image.repository = gke.gcr.io/gcp-compute-persistent-disk-csi-driver)
    
    - macro: known_aks_mount_in_privileged_containers
      condition:
        (k8s.ns.name = kube-system and container.image.repository = mcr.microsoft.com/oss/kubernetes-csi/azuredisk-csi
        or k8s.ns.name = system and container.image.repository = mcr.microsoft.com/oss/kubernetes-csi/secrets-store/driver)
    
    - macro: user_known_mount_in_privileged_containers
      condition: (never_true)
    
    - rule: Mount Launched in Privileged Container
      desc: Detect file system mount happened inside a privileged container which might lead to container escape.
      condition: >
        spawned_process and container
        and container.privileged=true
        and proc.name=mount
        and not mount_info
        and not known_gke_mount_in_privileged_containers
        and not known_aks_mount_in_privileged_containers
        and not user_known_mount_in_privileged_containers
      output: Mount was executed inside a privileged container (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid %container.info image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority: WARNING
      tags: [container, cis, filesystem, mitre_lateral_movement, T1611]
    
    - list: user_known_userfaultfd_processes
      items: []
    
    - rule: Unprivileged Delegation of Page Faults Handling to a Userspace Process
      desc: Detect a successful unprivileged userfaultfd syscall which might act as an attack primitive to exploit other bugs
      condition: >
        evt.type = userfaultfd and
        user.uid != 0 and
        (evt.rawres >= 0 or evt.res != -1) and
        not proc.name in (user_known_userfaultfd_processes)
      output: An userfaultfd syscall was successfully executed by an unprivileged user (user=%user.name user_loginuid=%user.loginuid process=%proc.name command=%proc.cmdline pid=%proc.pid %container.info image=%container.image.repository:%container.image.tag)
      priority: CRITICAL
      tags: [host, container, process, mitre_defense_evasion, TA0005]
    
    - list: ingress_remote_file_copy_binaries
      items: [wget]
    
    - macro: ingress_remote_file_copy_procs
      condition: (proc.name in (ingress_remote_file_copy_binaries))
    
    # Users should overwrite this macro to specify conditions under which a
    # Custom condition for use of ingress remote file copy tool in container
    - macro: user_known_ingress_remote_file_copy_activities
      condition: (never_true)
    
    - macro: curl_download
      condition: proc.name = curl and
                  (proc.cmdline contains " -o " or
                  proc.cmdline contains " --output " or
                  proc.cmdline contains " -O " or
                  proc.cmdline contains " --remote-name ")
    
    - rule: Launch Ingress Remote File Copy Tools in Container
      desc: Detect ingress remote file copy tools launched in container
      condition: >
        spawned_process and
        container and
        (ingress_remote_file_copy_procs or curl_download) and
        not user_known_ingress_remote_file_copy_activities
      output: >
        Ingress remote file copy tool launched in container (user=%user.name user_loginuid=%user.loginuid command=%proc.cmdline pid=%proc.pid parent_process=%proc.pname
        container_id=%container.id container_name=%container.name image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority: NOTICE
      tags: [container, network, process, mitre_command_and_control, TA0011]
    
    # This rule helps detect CVE-2021-4034:
    # A privilege escalation to root through memory corruption
    - rule: Polkit Local Privilege Escalation Vulnerability (CVE-2021-4034)
      desc: "This rule detects an attempt to exploit a privilege escalation vulnerability in Polkit's pkexec. By running specially crafted code, a local user can leverage this flaw to gain root privileges on a compromised system"
      condition:
        spawned_process and user.uid != 0 and proc.name=pkexec and proc.args = ''
      output:
        "Detect Polkit pkexec Local Privilege Escalation Exploit (CVE-2021-4034) (user=%user.loginname uid=%user.loginuid command=%proc.cmdline pid=%proc.pid args=%proc.args exe_flags=%evt.arg.flags)"
      priority: CRITICAL
      tags: [host, container, process, users, mitre_privilege_escalation, TA0004]
    
    
    - rule: Detect release_agent File Container Escapes
      desc: "This rule detect an attempt to exploit a container escape using release_agent file. By running a container with certains capabilities, a privileged user can modify release_agent file and escape from the container"
      condition:
        open_write and container and fd.name endswith release_agent and (user.uid=0 or thread.cap_effective contains CAP_DAC_OVERRIDE) and thread.cap_effective contains CAP_SYS_ADMIN
      output:
        "Detect an attempt to exploit a container escape using release_agent file (user=%user.name user_loginuid=%user.loginuid filename=%fd.name %container.info image=%container.image.repository:%container.image.tag cap_effective=%thread.cap_effective)"
      priority: CRITICAL
      tags: [container, process, mitre_privilege_escalation, mitre_lateral_movement, T1611]
    
    # Rule for detecting potential Log4Shell (CVE-2021-44228) exploitation
    # Note: Not compatible with Java 17+, which uses read() syscalls
    - macro: java_network_read
      condition: (evt.type=recvfrom and fd.type in (ipv4, ipv6) and proc.name=java)
    
    - rule: Java Process Class File Download
      desc: Detected Java process downloading a class file which could indicate a successful exploit of the log4shell Log4j vulnerability (CVE-2021-44228)
      condition: >
            java_network_read and evt.buffer bcontains cafebabe
      output: Java process class file download (user=%user.name user_loginname=%user.loginname user_loginuid=%user.loginuid event=%evt.type connection=%fd.name server_ip=%fd.sip server_port=%fd.sport proto=%fd.l4proto process=%proc.name command=%proc.cmdline pid=%proc.pid parent=%proc.pname buffer=%evt.buffer container_id=%container.id image=%container.image.repository)
      priority: CRITICAL
      enabled: false
      tags: [host, container, process, mitre_initial_access, T1190]
    
    - list: docker_binaries
      items: [docker, dockerd, containerd-shim, "runc:[1:CHILD]", pause, exe, docker-compose, docker-entrypoi, docker-runc-cur, docker-current, dockerd-current]
    
    - macro: docker_procs
      condition: proc.name in (docker_binaries)
    
    - rule: Modify Container Entrypoint
      desc: This rule detect an attempt to write on container entrypoint symlink (/proc/self/exe). Possible CVE-2019-5736 Container Breakout exploitation attempt.
      condition: >
        open_write and container and (fd.name=/proc/self/exe or fd.name startswith /proc/self/fd/) and not docker_procs and not proc.cmdline = "runc:[1:CHILD] init"
      enabled: false
      output: >
        Detect Potential Container Breakout Exploit (CVE-2019-5736) (user=%user.name process=%proc.name file=%fd.name cmdline=%proc.cmdline pid=%proc.pid %container.info)
      priority: WARNING
      tags: [container, filesystem, mitre_initial_access, T1611]
    
    - list: known_binaries_to_read_environment_variables_from_proc_files
      items: [scsi_id, argoexec]
    
    - rule: Read environment variable from /proc files
      desc: An attempt to read process environment variables from /proc files
      condition: >
        open_read and container and (fd.name glob /proc/*/environ)
        and not proc.name in (known_binaries_to_read_environment_variables_from_proc_files)
      output: >
        Environment variables were retrieved from /proc files (user=%user.name user_loginuid=%user.loginuid program=%proc.name
        command=%proc.cmdline pid=%proc.pid file=%fd.name parent=%proc.pname gparent=%proc.aname[2] ggparent=%proc.aname[3] gggparent=%proc.aname[4] container_id=%container.id image=%container.image.repository)
      priority: WARNING
      tags: [container, filesystem, process, mitre_credential_access, mitre_discovery, T1083]
    
    - list: known_ptrace_binaries
      items: []
    
    - macro: known_ptrace_procs
      condition: (proc.name in (known_ptrace_binaries))
    
    - macro: ptrace_attach_or_injection
      condition: >
        evt.type=ptrace and evt.dir=> and
        (evt.arg.request contains PTRACE_POKETEXT or
        evt.arg.request contains PTRACE_POKEDATA or
        evt.arg.request contains PTRACE_ATTACH or
        evt.arg.request contains PTRACE_SEIZE or
        evt.arg.request contains PTRACE_SETREGS)
    
    - rule: PTRACE attached to process
      desc: "This rule detects an attempt to inject code into a process using PTRACE."
      condition: ptrace_attach_or_injection and proc_name_exists and not known_ptrace_procs
      output: > 
        Detected ptrace PTRACE_ATTACH attempt (proc.cmdline=%proc.cmdline container=%container.info evt.type=%evt.type evt.arg.request=%evt.arg.request proc.pid=%proc.pid proc.cwd=%proc.cwd proc.ppid=%proc.ppid
        proc.pcmdline=%proc.pcmdline proc.sid=%proc.sid proc.exepath=%proc.exepath user.uid=%user.uid user.loginuid=%user.loginuid user.loginname=%user.loginname user.name=%user.name group.gid=%group.gid
        group.name=%group.name container.id=%container.id container.name=%container.name image=%container.image.repository)
      priority: WARNING
      tags: [host, container, process, mitre_execution, mitre_privilege_escalation, T1055.008]
      
    - rule: PTRACE anti-debug attempt
      desc: "Detect usage of the PTRACE system call with the PTRACE_TRACEME argument, indicating a program actively attempting to avoid debuggers attaching to the process. This behavior is typically indicative of malware activity."
      condition: evt.type=ptrace and evt.dir=> and evt.arg.request contains PTRACE_TRACEME and proc_name_exists
      output: Detected potential PTRACE_TRACEME anti-debug attempt (proc.cmdline=%proc.cmdline container=%container.info evt.type=%evt.type evt.arg.request=%evt.arg.request proc.pid=%proc.pid proc.cwd=%proc.cwd proc.ppid=%proc.ppid proc.pcmdline=%proc.pcmdline proc.sid=%proc.sid proc.exepath=%proc.exepath user.uid=%user.uid user.loginuid=%user.loginuid user.loginname=%user.loginname user.name=%user.name group.gid=%group.gid group.name=%group.name container.id=%container.id container.name=%container.name image=%container.image.repository)
      priority: NOTICE
      tags: [host, container, process, mitre_defense_evasion, T1622]
    
    - macro: private_aws_credentials
      condition: >
        (proc.args icontains "aws_access_key_id" or
        proc.args icontains "aws_secret_access_key" or
        proc.args icontains "aws_session_token" or
        proc.args icontains "accesskeyid" or
        proc.args icontains "secretaccesskey")
    
    - rule: Find AWS Credentials
      desc: Find or grep AWS credentials
      condition: >
        spawned_process and
        ((grep_commands and private_aws_credentials) or
        (proc.name = "find" and proc.args endswith ".aws/credentials"))
      output: Detected AWS credentials search activity (user.name=%user.name user.loginuid=%user.loginuid proc.cmdline=%proc.cmdline container.id=%container.id container_name=%container.name evt.type=%evt.type evt.res=%evt.res proc.pid=%proc.pid proc.cwd=%proc.cwd proc.ppid=%proc.ppid proc.pcmdline=%proc.pcmdline proc.sid=%proc.sid proc.exepath=%proc.exepath user.uid=%user.uid user.loginname=%user.loginname group.gid=%group.gid group.name=%group.name container.name=%container.name image=%container.image.repository:%container.image.tag exe_flags=%evt.arg.flags)
      priority: WARNING
      tags: [host, container, mitre_credential_access, process, aws, T1552]
    
    - rule: Execution from /dev/shm
      desc: This rule detects file execution from the /dev/shm directory, a common tactic for threat actors to stash their readable+writable+(sometimes)executable files.
      condition: >
        spawned_process and 
        (proc.exe startswith "/dev/shm/" or 
        (proc.cwd startswith "/dev/shm/" and proc.exe startswith "./" ) or 
        (shell_procs and proc.args startswith "-c /dev/shm") or 
        (shell_procs and proc.args startswith "-i /dev/shm") or 
        (shell_procs and proc.args startswith "/dev/shm") or 
        (proc.cwd startswith "/dev/shm/" and proc.args startswith "./" )) and 
        not container.image.repository in (falco_privileged_images, trusted_images)
      output: "File execution detected from /dev/shm (proc.cmdline=%proc.cmdline connection=%fd.name user.name=%user.name user.loginuid=%user.loginuid container.id=%container.id evt.type=%evt.type evt.res=%evt.res proc.pid=%proc.pid proc.cwd=%proc.cwd proc.ppid=%proc.ppid proc.pcmdline=%proc.pcmdline proc.sid=%proc.sid proc.exepath=%proc.exepath user.uid=%user.uid user.loginname=%user.loginname group.gid=%group.gid group.name=%group.name container.name=%container.name image=%container.image.repository exe_flags=%evt.arg.flags)"
      priority: WARNING
      tags: [host, container, mitre_execution, mitre_defense_evasion, T1036.005, T1059.004]
      
    # List of allowed container images that are known to execute binaries not part of their base image.
    # Users can use this list to better tune the rule below (i.e reducing false positives) by considering their workloads, 
    # since this requires application specific knowledge.
    - list: known_drop_and_execute_containers
      items: []
    
    - rule: Drop and execute new binary in container
      desc:
        Detect if an executable not belonging to the base image of a container is being executed.
        The drop and execute pattern can be observed very often after an attacker gained an initial foothold.
        is_exe_upper_layer filter field only applies for container runtimes that use overlayfs as union mount filesystem.
      condition: >
        spawned_process
        and container
        and proc.is_exe_upper_layer=true 
        and not container.image.repository in (known_drop_and_execute_containers)
      output: > 
        Executing binary not part of base image (user=%user.name user_loginuid=%user.loginuid user_uid=%user.uid comm=%proc.cmdline exe=%proc.exe container_id=%container.id
        image=%container.image.repository proc.name=%proc.name proc.sname=%proc.sname proc.pname=%proc.pname proc.aname[2]=%proc.aname[2] exe_flags=%evt.arg.flags
        proc.exe_ino=%proc.exe_ino proc.exe_ino.ctime=%proc.exe_ino.ctime proc.exe_ino.mtime=%proc.exe_ino.mtime proc.exe_ino.ctime_duration_proc_start=%proc.exe_ino.ctime_duration_proc_start
        proc.exepath=%proc.exepath proc.cwd=%proc.cwd proc.tty=%proc.tty container.start_ts=%container.start_ts proc.sid=%proc.sid proc.vpgid=%proc.vpgid evt.res=%evt.res)
      priority: CRITICAL 
      tags: [container, mitre_persistence, TA0003]
    
    # Application rules have moved to application_rules.yaml. Please look
    # there if you want to enable them by adding to
    # falco_rules.local.yaml.
    `,
)
