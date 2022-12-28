CREATE TABLE fooddo_entries(
  -- NOTE: rowid is the entry's unique id
  UserID BIGINT NOT NULL, -- Telegram's UserID of the entry's owner

  Date TEXT NOT NULL, -- In GO's format
  Type TINYINT CHECK(0 <= type AND type <= 4) NOT NULL, -- 0: Breakfast, 1: Lunch, 2: Merienda, 3: Dinner, 4: Extra
  Description TEXT NOT NULL,
  Skipped BOOL NOT NULL,

  Meat BOOL NOT NULL,
  Veggies BOOL NOT NULL,
  Fruit BOOL NOT NULL
);
