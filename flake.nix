{
  description = "Development shell & build Starlord";

  inputs = {
    nixpkgs.url     = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        isDarwin = builtins.match ".*-darwin.*" system != null;
      in {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; (
            [
              git
              kind
              kubectl
              k9s
              go-task
              go_1_23
              golint
            ]
          );
          shellHook = ''
            export IN_NIX_SHELL=1
            exec zsh
          '';
        };

      });
}