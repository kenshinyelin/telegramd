/*
 *  Copyright (c) 2017, https://github.com/nebulaim
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

package rpc

import (
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/grpc_util"
	"github.com/nebulaim/telegramd/mtproto"
	"golang.org/x/net/context"
	"github.com/nebulaim/telegramd/biz_model/dal/dao"
)

// auth.checkPhone#6fe51dfb phone_number:string = auth.CheckedPhone;
// tdesktop客户端会调用，android客户端未使用
func (s *AuthServiceImpl) AuthCheckPhone(ctx context.Context, request *mtproto.TLAuthCheckPhone) (*mtproto.Auth_CheckedPhone, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("AuthCheckPhone - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	phoneNumber, err := checkAndGetPhoneNumber(request.GetPhoneNumber())
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	usersDAO := dao.GetUsersDAO(dao.DB_SLAVE)
	usersDO := usersDAO.SelectByPhoneNumber(phoneNumber)
	glog.Infof("phoneNumber: %d, usersDO: {%v}", phoneNumber, usersDO)
	checkedPhone := mtproto.TLAuthCheckedPhone{Data2: &mtproto.Auth_CheckedPhone_Data{
		PhoneRegistered: mtproto.ToBool(usersDO != nil),
	}}

	glog.Infof("AuthCheckPhone - reply: %s\n", checkedPhone)
	return checkedPhone.To_Auth_CheckedPhone(), nil
}
