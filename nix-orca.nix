{
  lib,
  stdenv,
  fetchFromGitHub,
  buildGoModule,
}: let
  version = "0.0.4";
in
  buildGoModule {
    pname = "nix-orca";
    inherit version;

    # src = fetchFromGitHub {
    #   owner = "Verantor";
    #   repo = "nix-orca";
    #   rev = "v${version}";
    #   hash = "sha256-R4sqH4Nx0pQJc3UyuJ3dBd9PSEjyMmlaAlvforUJVH0=";
    # };
    src = ./.;

    vendorHash = "sha256-mN/QjzJ4eGfbW1H92cCKvC0wDhCR6IUes2HCZ5YBdPA=";
    # vendorHash = lib.fakeHash;
    meta = with lib; {
      homepage = "https://github.com/Verantor/nix-orca";
      description = "nix helper scripts";
      license = licenses.gpl3;
      maintainers = with maintainers; [Verantor];
    };
  }
