{ config, lib, pkgs, ... }: {
  # customize kernel version
  boot.kernelPackages = pkgs.linuxPackages_6_10;

  users = {
    users = {
      furry-profile = {
        isNormalUser = true;
        extraGroups = [ "wheel" "postgresql" ];
        password = "furry-profile";
        group = "furry-profile";
      };
    };

    groups.furry-profile = { };

  };

  virtualisation.vmVariant = {
    # following configuration is added only when building VM with build-vm
    virtualisation = {
      memorySize = 1024; # Use 2048MiB memory.
      cores = 1;
      graphics = false;
      forwardPorts = [
        { from = "host"; host.port = 5432; guest.port = 5432; }
      ];

    };
  };

  services = {
    openssh = {
      enable = true;
      settings.PasswordAuthentication = true;
    };

    postgresql = {
      enable = true;
      package = pkgs.postgresql_16;
      #extraPlugins = with config.services.postgresql.package.pkgs; [
      #  pg_repack
      #];
      # dataDir = "/home/furry-profile/datadir";
      # port = 5432;
      enableTCPIP = true;
      ensureDatabases = [ "furry-profile" ];
      ensureUsers = [{
        name = "furry-profile";
        ensureDBOwnership = true;
        ensureClauses.login = true;
      }];
      #settings = ''
      #  # CONFIG REFERENCE (with defaults)
      #  # https://github.com/postgres/postgres/blob/40ad4202513c72f5c1beeb03e26dfbc8890770c0/src/backend/utils/misc/postgresql.conf.sample
      #
      #  listen_addresses = '*'
      #  max_connections = 50
      #  jit = off # slows done things actually
      #
      #  logging_collector = on
      #  log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
      #  log_filename = 'postgresql-%Y-%m-%d.log'
      #  log_connections = on
      #  log_disconnections = on
      #  log_lock_waits = on
      #  log_temp_files = 0
      #  log_error_verbosity = default
      #  #log_min_messages = debug5
      #
      #
      #  # ssl = true
      #
      #  # those options should minimize WAL generation when inserting data
      #  # (WAL slows down everything, because it touches disk. But we still
      #  # cannot go with UNLOGGED tables, because those don't survive reboot)
      #  wal_level = minimal
      #  archive_mode = off
      #  max_wal_senders = 0
      #  wal_compression = on
      #  max_wal_size = 1GB
      #'';
      authentication = pkgs.lib.mkOverride 10 ''
        #type database  DBuser  auth-method
        local all       all     trust

        # TYPE  DATABASE        USER            ADDRESS                 METHOD
        host    all             all             all                     trust
      '';
    };
  };

  networking.firewall.allowedTCPPorts = [ 22 5432 ];

  environment.systemPackages = with pkgs; [
    htop
  ];

  # update this if needed
  system.stateVersion = "23.05";
}
