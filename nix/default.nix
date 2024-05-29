{ pkgs ? import <nixpkgs> {} }:
{
  nuv = pkgs.callPackage ./nuv.nix { };
}