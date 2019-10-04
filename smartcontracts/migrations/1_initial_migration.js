const Migrations = artifacts.require("Migrations");
const DecentralizedIdentifier = artifacts.require("DeviceDecentralizedIdentifier");
const IoTeXDID = artifacts.require("IoTeXDID");
const MockDeviceDID = artifacts.require("MockDeviceDID");
const CloudServiceAddr = "0x627306090abab3a6e1400e9345bc60c78a8bef57";

module.exports = function(deployer) {
  deployer.deploy(Migrations);
  deployer.deploy(DecentralizedIdentifier);
  deployer.deploy(IoTeXDID);
  deployer.deploy(MockDeviceDID, CloudServiceAddr);
};
