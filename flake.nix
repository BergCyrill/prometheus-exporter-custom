{
  description = "Project with Go and Task";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";  # or stable
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.go-task
          ];
          shellHook = ''
            export IN_NIX_SHELL=1
            exec zsh
          '';
        };
      }
    );
}