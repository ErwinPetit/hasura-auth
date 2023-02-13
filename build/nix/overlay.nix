(final: prev: rec {
  buildGoModule = final.buildGo119Module;

  go = prev.go_1_19.overrideAttrs
    (finalAttrs: previousAttrs: rec {
      version = "1.19.5";

      src = final.fetchurl {
        url = "https://go.dev/dl/go${version}.src.tar.gz";
        sha256 = "sha256-jkhujoWigfxc4/C+3FudLb9idtfbCyXT7ANPMT2gN18=";
      };

    });

  golangci-lint = prev.golangci-lint.override rec {
    buildGoModule = args: final.buildGoModule (args // rec {
      version = "1.50.1";
      src = final.fetchFromGitHub {
        owner = "golangci";
        repo = "golangci-lint";
        rev = "v${version}";
        sha256 = "sha256-7HoneQtKxjQVvaTdkjPeu+vJWVOZG3AOiRD87/Ntgn8=";
      };
      vendorHash = "sha256-6ttRd2E8Zsf/2StNYt6JSC64A57QIv6EbwAwJfhTDaY=";

      meta = with final.lib; args.meta // {
        broken = false;
      };
    });
  };

  golines = final.buildGoModule rec {
    name = "golines";
    version = "0.11.0";
    src = final.fetchFromGitHub {
      owner = "dbarrosop";
      repo = "golines";
      rev = "b7e767e781863a30bc5a74610a46cc29485fb9cb";
      sha256 = "sha256-pxFgPT6J0vxuWAWXZtuR06H9GoGuXTyg7ue+LFsRzOk=";
    };
    vendorSha256 = "sha256-rxYuzn4ezAxaeDhxd8qdOzt+CKYIh03A9zKNdzILq18=";

    meta = with final.lib; {
      description = "A golang formatter that fixes long lines";
      homepage = "https://github.com/segmentio/golines";
      maintainers = [ "nhost" ];
      platforms = platforms.linux ++ platforms.darwin;
    };
  };

  govulncheck = final.buildGoModule rec {
    name = "govulncheck";
    version = "latest";
    src = final.fetchFromGitHub {
      owner = "golang";
      repo = "vuln";
      rev = "dd534eeddf33556da7d61c8651a641b5e87be9d3";
      sha256 = "sha256-DeCZQWsdR4vjd+GRyRqz6vknqr6BEKPqBpSPnozN7+Q=";
    };
    vendorSha256 = "sha256-8022RUlhpr3hOEkjdfe4DXQ0K4G20EWBmyFltrqjF8M=";

    doCheck = false;

    meta = with final.lib; {
      description = "the database client and tools for the Go vulnerability database";
      homepage = "https://github.com/golang/vuln";
      maintainers = [ "nhost" ];
      platforms = platforms.linux ++ platforms.darwin;
    };
  };
})
