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

    src = fetchFromGitHub {
      owner = "Verantor";
      repo = "nix-orca";
      rev = "main";
      hash = "sha256-akrkZgqtChTkv1AWdrECjcMEgkktb0zi4kIIJ5pR5yc=";
    };
    # src = ./.;

    vendorHash = "sha256-PG6gCDZGLvWJh7iuaK60/yaGvshA2zracKlhmtQUtkU=";
    # vendorHash = lib.fakeHash;
    meta = with lib; {
      homepage = "https://github.com/Verantor/nix-orca";
      description = "nix helper scripts";
      license = licenses.gpl3;
      maintainers = with maintainers; [Verantor];
    };
  }
