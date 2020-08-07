# Deployment of DID Contract

## Install solc and ioctl
install solc according to https://docs.iotex.io/docs/ioctl.html#smart-contract

install ioctl according to https://docs.iotex.io/docs/ioctl.html#overview

## Deploy contract

1. First we need cd to smartcontracts/contracts directory.

2. We need to get the hexcode of "did:io:",that is "6469643a696f3a",io1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqd39ym7 is 0 address,that means the data contract will be created instead of using the existed.

3. Deploy contract using ioctl:

`ioctl contract deploy sol AddressBasedDIDManager AddressBasedDIDManager.sol DIDBase.sol DIDManagerBase.sol DIDStorage.sol IDID.sol ./ownership/Ownable.sol --with-arguments '{"_prefix":"6469643a696f3a","_dbAddr":"io1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqd39ym7"}' --endpoint api.testnet.iotex.one:80 --insecure -s acc1`

4. Get the deployed contract address from step 3 generated hash

`ioctl action hash fe007c9e3d551ff7a409812e62c9c5eac28c14145d308b841caeac16028e8bd1  --endpoint api.testnet.iotex.one:80 --insecure`
The contract address is io1w2xcw0meaklqxfj4ra98sz624uzxq5vcmvh483

5. Register DID to test the deployed contract if ok

`ioctl did register io1w2xcw0meaklqxfj4ra98sz624uzxq5vcmvh483 414efa99dfac6f4095d6954713fb0085268d400d6a05a8ae8a69b5b1c10b4bed uri -l 1000000 --endpoint api.testnet.iotex.one:80 --insecure -s acc2`

6. Get hash from the deployed contract,if the iotex address is io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms,we need to find it's ethereum address that is 0xda7e12ef57c236a06117c5e0d04a228e7181cf36,so the DID is did:io:0xda7e12ef57c236a06117c5e0d04a228e7181cf36

`ioctl did gethash io1w2xcw0meaklqxfj4ra98sz624uzxq5vcmvh483 did:io:0xda7e12ef57c236a06117c5e0d04a228e7181cf36 --endpoint api.testnet.iotex.one:80 --insecure`