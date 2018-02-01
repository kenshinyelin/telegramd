/*
 *  Copyright (c) 2018, https://github.com/nebulaim
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mtproto

import (
	"fmt"
)

type MessageBase interface {
	Encode() ([]byte, error)
	Decode(b []byte) error
}

//
//type CodecAble interface {
//	Encode() ([]byte, error)
//	Decode(dbuf *DecodeBuf) error
//}

func NewUnencryptedRawMessage() *UnencryptedRawMessage {
	return &UnencryptedRawMessage{
		AuthKeyId: 0,
	}
}

type UnencryptedRawMessage struct {
	// TODO(@benqi): reportAck and quickAck
	// NeedAck bool
	AuthKeyId int64
	MessageId int64
	MessageData []byte
}

func (m *UnencryptedRawMessage) Encode() ([]byte, error) {
	// 一次性分配
	x := NewEncodeBuf(20+len(m.MessageData))
	x.Long(0)
	m.MessageId = GenerateMessageId()
	x.Long(m.MessageId)
	x.Int(int32(len(m.MessageData)))
	x.Bytes(m.MessageData)
	return x.buf, nil
}

func (m *UnencryptedRawMessage) Decode(b []byte) error {
	dbuf := NewDecodeBuf(b)
	m.MessageId = dbuf.Long()
	messageLen := dbuf.Int()
	if int(messageLen) != dbuf.size-12 {
		return fmt.Errorf("invalid UnencryptedRawMessage len: %d (need %d)", messageLen, dbuf.size-12)
	}
	m.MessageData = dbuf.Bytes(int(messageLen))
	return dbuf.err
}

type EncryptedRawMessage struct {
	// TODO(@benqi): reportAck and quickAck
	// NeedAck bool
	AuthKeyId     int64
	MsgKey        []byte
	EncryptedData []byte
}

func NewEncryptedRawMessage(authKeyId int64) *EncryptedRawMessage {
	return &EncryptedRawMessage{
		AuthKeyId: authKeyId,
	}
}

func (m *EncryptedRawMessage) Encode() ([]byte, error) {
	// 一次性分配
	x := NewEncodeBuf(24+len(m.EncryptedData))
	x.Long(m.AuthKeyId)
	x.Bytes(m.MsgKey)
	x.Bytes(m.EncryptedData)
	return x.buf, nil
}

func (m *EncryptedRawMessage) Decode(b []byte) error {
	dbuf := NewDecodeBuf(b)
	m.MsgKey = dbuf.Bytes(16)
	m.EncryptedData = dbuf.Bytes(len(b)-16)
	return dbuf.err
}
