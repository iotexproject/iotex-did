const MockDeviceDID = artifacts.require('MockDeviceDID.sol');
const CloudServiceAddr = "0x886a06dbf45d6bed3d038488364c793b52221903"
const UUID = "WFFTUHU9BBC2PHJ3111A"
const Proof = "0x59e1106a8fdccdbafd34c766d6b28431e65b5467b900b05b205adcb8f19c55a149905101c20baa31ea122bb8eebef857a2b6ddef304c8fccbd69898e62a63a571b"

contract('mockdevicedid', function(accounts) {
    beforeEach(async function() {
        this.contract = await MockDeviceDID.new(CloudServiceAddr);
    });
    describe("create did", function () {
        it('success', async function() {
            let hash = web3.utils.fromAscii("hash");
            let uri = "s3://iotex-did/documents";
            let msg = 'I authorize' + accounts[0].toLowerCase() + ' to register ' + UUID;
            console.log(msg);
            let sig = await web3.eth.accounts.sign(msg, '0x1cd94a139f784fea91aee5a77b2519ab5852348a4525df12db2ef002d922e1e7');

            sig = sig.signature;
            console.log(sig);

            await this.contract.createDID(UUID, sig, hash, uri, {from: accounts[0]});
        })
    })
})