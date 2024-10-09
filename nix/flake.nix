# ref: https://gist.github.com/FlakM/0535b8aa7efec56906c5ab5e32580adf
{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      test-vm = nixpkgs.lib.nixosSystem {
        inherit system;
        modules = [
          ./vm-config.nix
        ];
      };
    in
    {
      # test is a hostname for our machine. This is optional if you don't need
      # to also expose the NixOS configuration for other purposes.
      nixosConfigurations.test = test-vm;

      # expose the build attribute directly
      # run with `nix build .#nixosConfigurations.test.config.system.build.vm`
      # `nix build .#vms.test`
      vms.test = test-vm.config.system.build.vm;
    };
}
