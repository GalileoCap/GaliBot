CREATE TABLE `users` (
  chatid bigint NOT NULL,
  name varchar(255) NOT NULL,
  permissions ENUM('admin', 'allow', 'block') NOT NULL,
  PRIMARY KEY (`chatid`)
);

INSERT INTO `users` (`chatid`, `name`, `permissions`)
VALUES (1129477471, 'Galileo', 'admin'),
       (5629879871, 'Sofia', 'admin');
