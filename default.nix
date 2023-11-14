{
  buildGoModule,
  pam,
  lib,
}:
buildGoModule {
  pname = "pam_k8s_sa";
  version = "unstable";

  src = lib.cleanSource ./.;

  vendorHash = "sha256-aXIV0unX+PWMNHED/vxzUcmLUWMONUhHCiwaZi8WMHI=";

  buildInputs = [
    pam
  ];

  buildPhase = ''
    runHook preBuild

    if [ -z "$enableParallelBuilding" ]; then
      export NIX_BUILD_CORES=1
    fi
    go build -buildmode=c-shared -o pam_k8s_sa.so -v -p $NIX_BUILD_CORES .

    runHook postBuild
  '';

  installPhase = ''
    runHook preInstall

    mkdir -p $out/lib/security
    cp pam_k8s_sa.so $out/lib/security

    runHook postInstall
  '';
}
