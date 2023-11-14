final: prev: let
  callPackage = final.darwin.apple_sdk_11_0.callPackage or final.callPackage;
in {
  pam_k8s_sa = callPackage ./default.nix {};
}
