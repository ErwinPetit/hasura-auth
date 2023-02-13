{
  description = "Nhost Hasura Auth";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/master";
    nix-filter.url = "github:numtide/nix-filter";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, nix-filter }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        localOverlay = import ./build/nix/overlay.nix;
        overlays = [ localOverlay ];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        go-src = nix-filter.lib.filter {
          root = ./.;
          include = with nix-filter.lib;[
            isDirectory
            (matchExt "go")
            (inDirectory "vendor")
            ./go.mod
            ./go.sum
          ];
        };

        nix-src = nix-filter.lib.filter {
          root = ./.;
          include = with nix-filter.lib;[
            (matchExt "nix")
          ];
        };

        goCheckDeps = with pkgs; [
          go
          clang
          golangci-lint
          richgo
          golines
          gofumpt
          govulncheck
          docker-compose
        ];

        jsCheckDeps = with pkgs; [
          nodejs-16_x
          nodePackages_latest.pnpm
        ];

        buildInputs = with pkgs; [
        ];

        nativeBuildInputs = with pkgs; [
          go
        ];

        name = "hasura-auth";
        version = nixpkgs.lib.fileContents ./VERSION;
        module = "github.com/nhost/hasura-auth";

        tags = [ ];

        ldflags = [
          "-X ${module}/controller.buildVersion=${version}"
        ];

      in
      {
        checks = {
          nixpkgs-fmt = pkgs.runCommand "check-nixpkgs-fmt"
            {
              nativeBuildInputs = with pkgs;
                [
                  nixpkgs-fmt
                ];
            }
            ''
              mkdir $out
              nixpkgs-fmt --check ${nix-src}
            '';

          goAuth = pkgs.runCommand "golang"
            {
              nativeBuildInputs = with pkgs; [
              ] ++ goCheckDeps ++ buildInputs ++ nativeBuildInputs;
            }
            ''
              export GOLANGCI_LINT_CACHE=$TMPDIR/.cache/golangci-lint
              export GOCACHE=$TMPDIR/.cache/go-build
              export GOMODCACHE="$TMPDIR/.cache/mod"
              export GOPATH="$TMPDIR/.cache/gopath"

              echo "➜ Source: ${go-src}"

              echo "➜ Running go generate ./... and checking sha1sum of all files"
              mkdir -p $TMPDIR/generate
              cd $TMPDIR/generate
              cp -r ${go-src}/* .
              chmod +w -R .

              go generate ./...
              find . -type f ! -path "./vendor/*" -print0 | xargs -0 sha1sum > $TMPDIR/sum
              cd ${go-src}
              sha1sum -c $TMPDIR/sum || (echo "❌ ERROR: go generate changed files" && exit 1)

              echo "➜ Running code formatters, if there are changes, fail"
              golines -l --base-formatter=gofumpt . | diff - /dev/null

              echo "➜ Checking for vulnerabilities"
              govulncheck ./...

              echo "➜ Running golangci-lint"
              golangci-lint run \
                --timeout 300s \
                ./...

              echo "➜ Running tests"
              richgo test \
                -tags="${pkgs.lib.strings.concatStringsSep " " tags}" \
                -ldflags="${pkgs.lib.strings.concatStringsSep " " ldflags}" \
                -v ./...

              mkdir $out
            '';

        };

        devShells = flake-utils.lib.flattenTree rec {
          build = pkgs.mkShell {
            buildInputs = with pkgs; [
            ] ++ buildInputs ++ nativeBuildInputs;
          };

          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              hey
            ] ++ jsCheckDeps ++ goCheckDeps ++ buildInputs ++ nativeBuildInputs;
          };
        };

        packages = flake-utils.lib.flattenTree rec {

          goAuth = pkgs.buildGoModule {
            inherit version ldflags buildInputs nativeBuildInputs;

            src = go-src;

            pname = name;

            doCheck = false;

            vendorHash = null;

            CGO_ENABLED = 0;

            meta = with pkgs.lib; {
              description = "Nhost Hasura Auth";
              homepage = "https://github.com/nhost/hasura-auth";
              maintainers = [ "nhost" ];
              platforms = platforms.linux ++ platforms.darwin;
            };
          };

          goAuth-docker-image = pkgs.dockerTools.buildLayeredImage {
            inherit name;
            tag = version;
            created = "now";

            contents = with pkgs; [
              pkgs.cacert
            ] ++ buildInputs;
            config = {
              Entrypoint = [
                "${self.packages.${system}.goAuth}/bin/hasura-auth"
              ];
            };
          };


        };


      }



    );


}
