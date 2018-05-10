CREATE TABLE `user` (
  `id` varchar(128) NOT NULL COMMENT 'id',
  `username` varchar(45) NOT NULL COMMENT '用户名\n',
  `userid` varchar(60) NOT NULL COMMENT '登录Id',
  `email` varchar(60) NOT NULL COMMENT '电子邮件',
  `mobileno` varchar(16) NOT NULL COMMENT '手机',
  `address` varchar(45) DEFAULT NULL COMMENT '地址',
  `sex` tinyint(1) NOT NULL COMMENT '姓别\n0：男 1：女',
  `introduce` varchar(200) DEFAULT NULL COMMENT '介绍',
  `avatar` varchar(150) DEFAULT NULL COMMENT '头像地址url',
  `bgimgurl` varchar(150) DEFAULT NULL COMMENT '背景图url',
  `job` varchar(45) DEFAULT NULL COMMENT '职业',
  `city` varchar(45) DEFAULT NULL COMMENT '城市',
  `country` varchar(45) DEFAULT NULL COMMENT '国家',
  `createtime` date NOT NULL COMMENT '创建时间',
  `updatetime` date DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;




CREATE TABLE `group` (
  `id` varchar(128) NOT NULL COMMENT 'id',
  `name` varchar(45) NOT NULL COMMENT '名称',
  `describe` varchar(150) NOT NULL COMMENT '描述',
  `createtime` date NOT NULL COMMENT '创建时间',
  `createuserid` varchar(128) NOT NULL COMMENT '创建人用户id',
  `state` varchar(2) NOT NULL COMMENT '状态 0:正常 ； 1:冻结',
  `updatetime` date DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `groupid_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;






CREATE TABLE `group_member` (
  `groupid` varchar(128) NOT NULL COMMENT '群组id\n',
  `userid` varchar(128) NOT NULL COMMENT '用户Id',
  `nickname` varchar(120) DEFAULT NULL COMMENT '用户在群组内的别称',
  `card` varchar(180) DEFAULT NULL COMMENT '用户在群组内的名片\n',
  `jointime` date NOT NULL COMMENT '加入群的时间',
  `state` varchar(2) DEFAULT NULL COMMENT '在群内的状态\n0：正常  1：禁言',
  `onlinetime` double NOT NULL COMMENT '总在线时间 单位小时',
  `level` varchar(45) DEFAULT NULL COMMENT '等级',
  `updatetime` date DEFAULT NULL COMMENT '更新时间',
  `manageflag` varchar(1) NOT NULL COMMENT '管理标识\n0：普通成员 1：管理员  2:超级管理员',
  PRIMARY KEY (`groupid`,`userid`),
  UNIQUE KEY `groupid_UNIQUE` (`groupid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



CREATE TABLE `user_manage_jnl` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '流水序号',
  `userid` varchar(128) NOT NULL COMMENT '用户Id\n',
  `transcode` varchar(45) NOT NULL COMMENT '交易码',
  `transdate` date NOT NULL COMMENT '交易时间',
  `reqdata` json NOT NULL COMMENT '请求数据 json格式',
  `state` varchar(45) NOT NULL COMMENT '状态\n0：成功 1：失败',
  `updatetime` varchar(45) DEFAULT NULL COMMENT '更新时间',
  `errmsg` varchar(45) DEFAULT NULL COMMENT '错误信息\n',
  PRIMARY KEY (`id`,`userid`,`transdate`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;




CREATE TABLE `group_manage_jnl` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '流水号\n',
  `groupid` varchar(128) NOT NULL COMMENT '群组id\n',
  `transcode` varchar(45) DEFAULT NULL COMMENT '交易码',
  `transdate` date NOT NULL COMMENT '交易时间\n',
  `reqdata` json NOT NULL COMMENT '请求数据',
  `state` varchar(2) DEFAULT NULL COMMENT '状态\n',
  `errmsg` varchar(45) DEFAULT NULL COMMENT '错误信息',
  `updatetime` varchar(45) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`,`groupid`,`transdate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `group_member_manage_jnl` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '序号',
  `groupid` varchar(128) NOT NULL COMMENT '群组id',
  `userid` varchar(45) NOT NULL COMMENT '用户id',
  `transcode` varchar(45) NOT NULL COMMENT '交易码\n',
  `reqdata` json NOT NULL COMMENT '请求数据',
  `transdate` date NOT NULL,
  `state` varchar(2) DEFAULT NULL COMMENT '状态\n0：成功  1：失败',
  `errmsg` varchar(100) DEFAULT NULL COMMENT '错误信息\n',
  `updatetime` date DEFAULT NULL,
  PRIMARY KEY (`id`,`groupid`,`userid`,`transdate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;