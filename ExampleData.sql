INSERT INTO restaurants (id, name) VALUES
  (1, 'Pilars Place'),
  (2, 'Mikes Bistro'),
  (3, 'Mollys'),
  (4, 'Murphys'),
  (5, 'Sushi Ya'),
  (6, 'Umplebees Bakery');

INSERT INTO users (restaurant_id, username, password, pass_salt) VALUES
  (1, 'michaelp', 'fakepass', 'salt'),
  (3, 'georgeb', 'fakepass', 'salt'),
  (1, 'donw', 'fakepass', 'salt'),
  (2, 'hugom', 'fakepass', 'salt'),
  (6, 'johnw', 'fakepass', 'salt');


INSERT INTO buzzers (restaurant_id, buzzer_name, is_active) VALUES
        (3, 'loud-turtke-1043', 'FALSE'),
        (3, 'smart-fox-2066', 'FALSE'),
        (3, 'fast-dog-3212', 'TRUE');

INSERT INTO active_parties (restaurant_id, party_name, party_size, time_created, phone_ahead, wait_time_expected, wait_time_calculated, is_table_ready, buzzer_id) VALUES
        (1, 'Josh', 11, '2016-10-30 08:43:24', 'FALSE', 27, 25, 'FALSE', 1),
        (3, 'Mickey', 4, '2016-10-30 11:30:00', 'FALSE', 10, 15, 'FALSE', 2),
        (1, 'Kathy', 3, '2016-10-30 10:13:54', 'TRUE', 20, 25, 'FALSE', 3),
        (9, 'Joe', 5, '2016-10-30 10:23:54', 'FALSE', 40, 35, 'FALSE', 4),
        (9, 'Bob', 2, '2016-10-30 09:23:54', 'FALSE', 10, 15, 'FALSE', 5),
        (9, 'Mark', 3, '2016-10-30 09:43:54', 'TRUE', 20, 25, 'FALSE', 6);

INSERT INTO historical_parties (restaurant_id, party_name, party_size, time_created, time_seated, wait_time_expected, wait_time_calculated) VALUES
    (1, 'Josh', 11, '2016-10-30 08:43:24', '2016-10-30 09:43:24', 27, 25),
    (1, 'Larson', 3, '2016-10-30 10:16:54','2016-10-30 10:59:24', 20, 25),
    (1, 'Dude', 3, '2016-10-30 10:17:54', '2016-10-30 11:43:24',20, 29),
    (1, 'Bro', 3, '2016-10-30 10:18:54', '2016-10-30 10:30:54',20, 25),
    (1, 'Bistroman', 3, '2016-10-30 10:19:54', '2016-10-30 10:44:54', 20, 25),
    (1, 'Sausageking', 3, '2016-10-30 10:20:54', '2016-10-30 10:40:54', 20, 25),
    (1, 'Mickey', 4, '2016-10-30 11:21:00', '2016-10-30 11:34:00', 10, 15);

