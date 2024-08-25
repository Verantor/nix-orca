{
  lib,
  stdenv,
  fetchFromGitHub,
  buildGoModule,
}: let
  version = "0.0.3";
in
  buildGoModule {
    pname = "nix-orca";
    inherit version;

    src = fetchFromGitHub {
      owner = "Verantor";
      repo = "nix-orca";
      rev = "v${version}";
      # hash = "sha256-te6fxQ99sJEFCdm29Qjja8XAXz8NS0kji8CyV2Cnfyc=";
    };

    vendorHash = "sha256-mN/QjzJ4eGfbW1H92cCKvC0wDhCR6IUes2HCZ5YBdPA";
    meta = with lib; {
      homepage = "https://github.com/Verantor/nix-orca";
      description = "nix helper scripts";
      license = licenses.gpl3;
      maintainers = with maintainers; [Verantor];
    };
  }
