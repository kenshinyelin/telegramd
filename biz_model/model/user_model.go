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

package model

import (
	"sync"
	"github.com/nebulaim/telegramd/biz_model/dal/dao"
	"github.com/nebulaim/telegramd/mtproto"
	// "github.com/nebulaim/telegramd/base/base"
	// "github.com/golang/glog"
	// "github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/biz_model/dal/dataobject"
	"github.com/nebulaim/telegramd/baselib/base"
)

type userModel struct {
	// usersDAO *dao.UsersDAO
}

var (
	userInstance *userModel
	userInstanceOnce sync.Once
)

func GetUserModel() *userModel {
	userInstanceOnce.Do(func() {
		userInstance = &userModel{}
	})
	return userInstance
}

// func (m *userModel) User(userId int32) (user* mtproto.TLUser) {

func (m *userModel) GetUser(userId int32) (user* mtproto.TLUser) {
	usersDAO := dao.GetUsersDAO(dao.DB_SLAVE)
	userDO := usersDAO.SelectById(userId)

	if userDO != nil {
		// TODO(@benqi): fill bot, photo, about...
		user = &mtproto.TLUser{ Data2: &mtproto.User_Data{
			// user.Self由业务层进行判断
			// user.Self = true
			Id: userDO.Id,
			AccessHash: userDO.AccessHash,
			FirstName: userDO.FirstName,
			LastName: userDO.LastName,
			Username: userDO.Username,
			Phone: userDO.Phone,
		}}
	}
	return
}

func (m *userModel) GetUsersBySelfAndIDList(selfUserId int32, userIdList []int32) (users []*mtproto.User) {
	if len(userIdList) == 0 {
		users = []*mtproto.User{}
	} else {
		// usersDAO := dao.GetUsersDAO(dao.DB_SLAVE)
		userDOList := dao.GetUsersDAO(dao.DB_SLAVE).SelectUsersByIdList(userIdList)
		users = make([]*mtproto.User, 0, len(userDOList))
		for _, userDO := range userDOList {
			// TODO(@benqi): fill bot, photo, about...
			user := &mtproto.TLUser{Data2: &mtproto.User_Data{
				Self:          selfUserId == userDO.Id,
				Id:            userDO.Id,
				AccessHash:    userDO.AccessHash,
				FirstName:     userDO.FirstName,
				LastName:      userDO.LastName,
				Username:      userDO.Username,
				Phone:         userDO.Phone,
				Contact:       true,
				MutualContact: true,
			}}
			users = append(users, user.To_User())
		}
	}
	return
}

func (m *userModel) GetUserList(userIdList []int32) (users []*mtproto.TLUser) {
	usersDAO := dao.GetUsersDAO(dao.DB_SLAVE)

	userDOList := usersDAO.SelectUsersByIdList(userIdList)
	users = []*mtproto.TLUser{}
	for _, userDO := range userDOList {
		// TODO(@benqi): fill bot, photo, about...
		user := &mtproto.TLUser{Data2: &mtproto.User_Data{
			// user.Self由业务层进行判断
			// user.Self = true
			Id:         userDO.Id,
			AccessHash: userDO.AccessHash,
			FirstName:  userDO.FirstName,
			LastName:   userDO.LastName,
			Username:   userDO.Username,
			Phone:      userDO.Phone,
		}}

		users = append(users, user)
	}

	// glog.Infof("SelectUsersByIdList(%s) - %s", base.JoinInt32List(userIdList, ","), logger.JsonDebugData(users))
	return
}

func (m *userModel) GetUserFull(userId int32) (userFull *mtproto.TLUserFull) {
	//TODO(@benqi): 等Link和NotifySettings实现后再来完善
	//fullUser := &mtproto.TLUserFull{}
	//fullUser.PhoneCallsAvailable = true
	//fullUser.PhoneCallsPrivate = true
	//fullUser.About = "@Benqi"
	//fullUser.CommonChatsCount = 0
	//
	//switch request.Id.Payload.(type) {
	//case *mtproto.InputUser_InputUserSelf:
	//	// User
	//	userDO, _ := s.UsersDAO.SelectById(2)
	//	user := &mtproto.TLUser{}
	//	user.Self = true
	//	user.Contact = false
	//	user.Id = userDO.Id
	//	user.FirstName = userDO.FirstName
	//	user.LastName = userDO.LastName
	//	user.Username = userDO.Username
	//	user.AccessHash = userDO.AccessHash
	//	user.Phone = userDO.Phone
	//	fullUser.User = mtproto.MakeUser(user)
	//
	//	// Link
	//	link := &mtproto.TLContactsLink{}
	//	link.MyLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.ForeignLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.User = mtproto.MakeUser(user)
	//	fullUser.Link = mtproto.MakeContacts_Link(link)
	//case *mtproto.InputUser_InputUser:
	//	inputUser := request.Id.Payload.(*mtproto.InputUser_InputUser).InputUser
	//	// User
	//	userDO, _ := s.UsersDAO.SelectById(inputUser.UserId)
	//	user := &mtproto.TLUser{}
	//	user.Self = false
	//	user.Contact = true
	//	user.Id = userDO.Id
	//	user.FirstName = userDO.FirstName
	//	user.LastName = userDO.LastName
	//	user.Username = userDO.Username
	//	user.AccessHash = userDO.AccessHash
	//	user.Phone = userDO.Phone
	//	fullUser.User = mtproto.MakeUser(user)
	//
	//	// Link
	//	link := &mtproto.TLContactsLink{}
	//	link.MyLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.ForeignLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.User = mtproto.MakeUser(user)
	//	fullUser.Link = mtproto.MakeContacts_Link(link)
	//case *mtproto.InputUser_InputUserEmpty:
	//	// TODO(@benqi): BAD_REQUEST: 400
	//}
	//
	//// NotifySettings
	//peerNotifySettings := &mtproto.TLPeerNotifySettings{}
	//peerNotifySettings.ShowPreviews = true
	//peerNotifySettings.MuteUntil = 0
	//peerNotifySettings.Sound = "default"
	//fullUser.NotifySettings = mtproto.MakePeerNotifySettings(peerNotifySettings)
	return nil
}

func (m *userModel) UpdateUserStatus(userId int32, lastSeenAt int64) {
	presencesDAO := dao.GetUserPresencesDAO(dao.DB_MASTER)
	// now := time.Now().Unix()
	rows := presencesDAO.UpdateLastSeen(lastSeenAt, 0, userId)
	if rows == 0 {
		do := &dataobject.UserPresencesDO{
			UserId: userId,
			LastSeenAt: lastSeenAt,
			LastSeenAuthKeyId: 0,
			LastSeenIp: "",
			CreatedAt: base.NowFormatYMDHMS(),
		}
		presencesDAO.Insert(do)
	}
}
