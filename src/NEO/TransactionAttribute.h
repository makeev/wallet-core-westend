// Copyright © 2017-2020 Trust Wallet.
//
// This file is part of Trust. The full Trust copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

#pragma once

#include "TransactionAttributeUsage.h"
#include "ISerializable.h"
#include "Serializable.h"
#include "Data.h"

namespace TW::NEO {

class TransactionAttribute : public Serializable {
  public:
    TransactionAttributeUsage usage = TAU_ContractHash;
    Data _data;

    virtual ~TransactionAttribute() {}

    int64_t size() const override {
        return 1 + _data.size();
    }

    void deserialize(const Data& data, int initial_pos = 0) override {
        if (static_cast<int>(data.size()) < initial_pos + 1) {
            throw std::invalid_argument("Invalid data for deserialization");
        }
        usage = (TransactionAttributeUsage) data[initial_pos];
        if (usage == TransactionAttributeUsage::TAU_ContractHash || usage == TransactionAttributeUsage::TAU_Vote ||
            (usage >= TransactionAttributeUsage::TAU_Hash1 && usage <= TransactionAttributeUsage::TAU_Hash15)) {
            _data = readBytes(data, 32, initial_pos + 1);
        } else if (usage == TransactionAttributeUsage::TAU_ECDH02 ||
                    usage == TransactionAttributeUsage::TAU_ECDH03) {
            _data = readBytes(data, 32, initial_pos + 1);
        } else if (usage == TransactionAttributeUsage::TAU_Script) {
            _data = readBytes(data, 20, initial_pos + 1);
        } else if (usage == TransactionAttributeUsage::TAU_DescriptionUrl) {
            _data = readBytes(data, 1, initial_pos + 1);
        } else if (usage == TransactionAttributeUsage::TAU_Description ||
                    usage >= TransactionAttributeUsage::TAU_Remark) {
            _data = readBytes(data, int(data.size()) - 1 - initial_pos, initial_pos + 1);
        } else {
            throw std::invalid_argument("TransactionAttribute Deserialize FormatException");
        }
    }

    Data serialize() const override {
        return concat(Data({static_cast<byte>(usage)}), _data);
    }

    bool operator==(const TransactionAttribute &other) const {
        return this->usage == other.usage
            && _data.size() == other._data.size()
            && _data == other._data;
    }
};

}
