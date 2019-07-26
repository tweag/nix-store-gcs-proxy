with import <nixpkgs> {};
buildGoPackage {
  name = "nix-store-gcs-proxy";
  goPackagePath = "github.com/tweag/nix-store-gcs-proxy";
  src = ./.;
  goDeps = ./deps.nix;
}
