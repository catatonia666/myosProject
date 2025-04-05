CREATE TABLE starting_blocks (
    id SERIAL PRIMARY KEY,
    story_title TEXT NOT NULL,
    user_id INT NOT NULL,
    privacy BOOLEAN NOT NULL DEFAULT FALSE,
    content TEXT NOT NULL,
    options JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE common_blocks (
    id SERIAL PRIMARY KEY,
    story_id INT NOT NULL REFERENCES starting_blocks(id) ON DELETE CASCADE,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    options JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO starting_blocks (story_title, id, user_id, privacy, content, options)
VALUES (
    'A Story About Learning', 
    1, 
    0, 
    false, 
    'Hello dear travaller! It is very nice to see you in my log cabin. This is where I create my stories, and here you can create yours! I am eager to aswer your questions, if you have them.', 
    '[{"1": "How can I see the stories of other creators?"}, {"2": "How can I create my own story?"}]'
);

SELECT setval('public.starting_blocks_id_seq', (SELECT MAX(id) FROM starting_blocks));

INSERT INTO common_blocks (id, story_id, user_id, content, options)
VALUES (
    1, 
    1, 
    0, 
    'Great question endeed. There are a lot of talanted creators around here, and all of them wants to share their stories with each other. Luckily it is very simple to accomplish. All you need to do is.. Well, actually you dont need to do anything. Just visit the home page where all non-private stories are waiting for you', 
    '[{"3": "Can I edit someones story?"}, {"4": "I see, thank you!"}]'
);

INSERT INTO common_blocks (id, story_id, user_id, content, options)
VALUES (
    2, 
    1, 
    0, 
    'To create your own stories you should first to sign up and then log in.  Then just press New story button and create! But it is higly recommended to read About page first, because some tools can be a bit.. tricky.', 
    '[{"5": "Can I add some pictures?"}, {"5": "Is it possible to add any pictures in my story?"}]'
);

INSERT INTO common_blocks (id, story_id, user_id, content, options)
VALUES (
    3, 
    1, 
    0, 
    'No, you can not edit anyones stories, but feel free to edit your own. Just log in and feel free to correct anything you want.', 
    '[{"4": "Well, thank you then!"}]'
);

INSERT INTO common_blocks (id, story_id, user_id, content, options)
VALUES (
    4, 
    1, 
    0, 
    'You are welcome!', 
    '[]'
);

INSERT INTO common_blocks (id, story_id, user_id, content, options)
VALUES (
    5, 
    1, 
    0, 
    'No, unfortunatly it is not possible yet, but I think this functionality will appear any time soon!', 
    '[{"4": "Well, thank you then!"}]'
);

 SELECT setval('public.common_blocks_id_seq', (SELECT MAX(id) FROM common_blocks));