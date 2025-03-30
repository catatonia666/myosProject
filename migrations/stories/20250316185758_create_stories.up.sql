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
    'Secrets of the Omnissiah', 
    1, 
    1, 
    false, 
    '++ INCOMING DATA TRANSMISSION ++\n\nYou stand before a towering figure of metal and crimson robes. A tech-priest of the Adeptus Mechanicus stares at you, mechadendrites twitching, eyes glowing a dim red. "State your query, fleshling," his voice crackles with binharic undertones.', 
    '[{"2": "What is the Omnissiah?"}, {"3": "Why do you replace your flesh with metal?"}]'
);

SELECT setval('public.starting_blocks_id_seq', (SELECT MAX(id) FROM starting_blocks));

INSERT INTO common_blocks (id, story_id, user_id, content, options)
VALUES (
    2, 
    1, 
    1, 
    '"The Omnissiah is the great Machine God, the source of all knowledge and divinity in the universe. It is through His divine algorithms that we comprehend the sacred laws of technology. The ignorant may call Him the Emperor, but true understanding is reserved for those who seek enlightenment through the cog and circuit."\n\nThe tech-priest’s servo-skull hovers near you, scanning for impurities.', 
    '[{"4": "What is the Quest for Knowledge?"}, {"5": "What happens to those who deny the Omnissiah?"}]'
);

INSERT INTO common_blocks (id, story_id, user_id, content, options)
VALUES (
    3, 
    1, 
    1, 
    '"Flesh is weak. Flesh decays. But the machine endures. Through augmentation, we ascend beyond the frailties of mortality. The more of our feeble bodies we replace, the closer we become to the divine perfection of the Omnissiah."\n\nThe tech-priest’s metallic fingers flex with a quiet hiss of hydraulics.', 
    '[{"6": "Do you feel anything at all?"}, {"7": "How do you communicate with machines?"}]'
);

INSERT INTO common_blocks (id, story_id, user_id, content, options) VALUES (4, 1, 1, '', '{}');
INSERT INTO common_blocks (id, story_id, user_id, content, options) VALUES (5, 1, 1, '', '{}');
INSERT INTO common_blocks (id, story_id, user_id, content, options) VALUES (6, 1, 1, '', '{}');
INSERT INTO common_blocks (id, story_id, user_id, content, options) VALUES (7, 1, 1, '', '{}');

 SELECT setval('public.common_blocks_id_seq', (SELECT MAX(id) FROM common_blocks));