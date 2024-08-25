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
      owner = "arsham";
      repo = "figurine";
      rev = "v${version}";
      hash = "sha256-1q6Y7oEntd823nWosMcKXi6c3iWsBTxPnSH4tR6+XYs=";
    };

    vendorSha256 = "sha256-mLdAaYkQH2RHcZft27rDW1AoFCWKiUZhh2F0DpqZELw=";

    meta = with lib; {
      homepage = "https://github.com/arsham/figurine";
      description = "Print your name in style";
      license = licenses.asl20;
      maintainers = with maintainers; [ironicbadger];
    };
  }
