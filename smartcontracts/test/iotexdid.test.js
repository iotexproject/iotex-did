const IoTeXDID = artifacts.require('IoTeXDID.sol');

contract('iotexdid', function(accounts) {
    beforeEach(async function() {
        this.contract = await IoTeXDID.new();
    });
    describe("create did", function () {
        it('success', async function() {
            let hash = web3.utils.fromAscii("hash");
            let uri = "s3://iotex-did/documents";
            await this.contract.createDID(accounts[0], hash, uri, {from: accounts[0]});
        })
    })
})