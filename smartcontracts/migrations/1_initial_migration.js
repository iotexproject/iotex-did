const Migrations = artifacts.require("Migrations");
const DecentralizedIdentifier = artifacts.require("DecentralizedIdentifier");
const IoTeXDID = artifacts.require("IoTeXDID");

module.exports = function(deployer) {
  deployer.deploy(Migrations);
  deployer.deploy(DecentralizedIdentifier);
  deployer.deploy(IoTeXDID);
};
