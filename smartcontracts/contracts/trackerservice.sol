pragma solidity ^0.4.24;

import './ownership/Whitelist.sol';

contract TrackerService is Whitelist {
    uint256 unitPrice;
    uint256 maximumAccessPeriod;

    event UpdateUnitPrice(uint256 unitPrice);
    event UpdateMaximumAccessPeriod(uint256 maximumAccessPeriod);
    event RequestDataPointAccess(string indexed publicKey, uint256 expirationTime);
    event ShareAccessToken(string associationToken, string indexed tokenUserPublicKey);

    constructor(uint256 _unitPrice, uint256 _maximumAccessPeriod) public {
        addAddressToWhitelist(msg.sender);
        unitPrice = _unitPrice;
        maximumAccessPeriod = _maximumAccessPeriod;
    }

    function updateUnitPrice(uint256 _unitPrice) public onlyWhitelisted {
        unitPrice = _unitPrice;
        emit UpdateUnitPrice(_unitPrice);
    }

    function updateMaximumAccessPeriod(uint256 _maximumAccessPeriod) public onlyWhitelisted {
        maximumAccessPeriod = _maximumAccessPeriod;
        emit UpdateMaximumAccessPeriod(_maximumAccessPeriod);
    }

    function requestDataPointAccess(string publicKey, uint256 expirationTime) public payable {
        if (unitPrice != 0) {
            require(msg.value % unitPrice == 0, "invalid payment amount");
        }
        require(expirationTime > now && expirationTime - now <= maximumAccessPeriod, "invalid expiration time");

        emit RequestDataPointAccess(publicKey, expirationTime);
    }

    function shareAccessToken(string associationToken, string tokenUserPublicKey) public onlyWhitelisted {
        emit ShareAccessToken(associationToken, tokenUserPublicKey);
    }
}