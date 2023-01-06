CREATE TABLE IF NOT EXISTS users(
  id BIGINT UNIQUE NOT NULL, -- Telegram's UserID

  firstname VARCHAR(64) NOT NULL DEFAULT "",
  lastname VARCHAR(64) NOT NULL DEFAULT "",
  username VARCHAR(32) NOT NULL DEFAULT "", 

  permissions TINYINT CHECK(0 <= permissions AND permissions <= 3) NOT NULL -- 0: admin, 1: tester, 2: user, 3: blocked
);

INSERT OR REPLACE INTO users(id, firstname, lastname, username, permissions)
VALUES
  (1129477471, "Galileo", "Cappella", "galileocap", 0),
  (5629879871, "Sofia", "B", "", 1)
