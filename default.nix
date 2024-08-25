{
  lib,
  stdenv,
  fetchFromGitHub,
  buildGoModule,
}: let
  version = "0.0.1";
in
  buildGoModule {
    pname = "lynx";
    inherit version;

    src = fetchFromGitHub {
      owner = "Verantor";
      repo = "nix-orca";
      rev = "v${version}";
      # hash = "sha256-1q6Y7oEntd823nWosMcKXi6c3iWsBTxPnSH4tR6+XYs=";
    };

    vendorHash = lib.fakeHash;
    meta = with lib; {
      homepage = "https://github.com/Verantor/nix-orca";
      description = "nix helper scripts";
      license = licenses.gpl3;
      maintainers = with maintainers; [Verantor];
    };
  }
