{
  description = "Monorepo plugin for buildkite";

  inputs = {
    # Nix Inputs
    nixpkgs.url = github:nixos/nixpkgs/nixpkgs-unstable;
    flake-utils.url = github:numtide/flake-utils;
  };

  outputs = {
    flake-utils,
    nixpkgs,
    self
  }:
    with flake-utils.lib;
      eachSystem
      [
        "aarch64-linux"
        "aarch64-darwin"
        "x86_64-darwin"
        "x86_64-linux"
      ]
      (system: let
        pkgs = import nixpkgs {
          inherit system;
        };
        version = "2.5.1";
        monorepo-diff = pkgs.buildGo118Module ({
          inherit version;
          src = ./.;
          vendorSha256 = "sha256-iEUbfIRJtyvXoE5VHra+07SIXmGSpWOHhVDI61vh0Ck=";
          pname = "monorepo-diff";
          ldflags = ''
            -X main.Version=${version}
            -X main.BuildSha=${if (self ? rev) then self.rev else throw "You must commit git changes to run this binary"}
          '';
        });
      in {
        apps = {
          default = mkApp {
            name = "monorepo-diff";
            drv = monorepo-diff;
          };
        };
      });
}

