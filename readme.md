```mysql
CREATE TABLE `dtb_permission_rule` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `title` longtext,
  `name` longtext,
  `parent_id` int(11) DEFAULT '0',
  `status` int(11) DEFAULT '0',
  `level` int(11) DEFAULT NULL COMMENT '级别。1模块,2控制器,3操作',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;

CREATE TABLE `dtb_permission_role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `title` longtext COMMENT '名称',
  `memo` longtext COMMENT '组描述',
  `status` int(11) DEFAULT '0' COMMENT '启用状态',
  `rules` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

CREATE TABLE `dtb_member_role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `member_id` int(11) DEFAULT NULL,
  `permission_role_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
```


```mysql
INSERT INTO `dtb_permission_rule` (`id`, `title`, `name`, `parent_id`, `status`, `level`) VALUES ('1', '首页', 'dashboard', '0', '1', '1'),
('2', '仪表盘', 'dashboard', '1', '1', '2'),
('3', '列表', 'index', '2', '1', '3'),
('4', '文件管理', 'content', '0', '1', '1'),
('5', '文件管理', 'file', '4', '1', '2'),
('6', '上传', 'upload', '5', '1', '3'),
('7', '下载', 'download', '5', '1', '3'),
('8', '预览', 'view', '5', '1', '3'),
('9', '列表', 'index', '5', '1', '3'),
('10', '删除', 'delete', '5', '1', '3'),
('11', '创建文件夹', 'createfolder', '5', '1', '3'),
('12', '创建文件', 'touchfile', '5', '1', '3'),
('13', '平台管理', 'platform', '0', '1', '1'),
('14', '系统管理', 'system', '13', '1', '1'),
('15', '系统环境', 'environment', '14', '1', '2'),
('16', '列表', 'index', '15', '1', '3'),
('17', '重启', 'restart', '15', '1', '3');
```