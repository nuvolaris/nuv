{ pkgs ? import <nixpkgs> {} }:
let
  nuv = pkgs.callPackage ./nuv.nix { }; 
in
pkgs.mkShell {
  buildInputs = [
    nuv
  ];
}
