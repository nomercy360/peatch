CREATE TABLE users
(
    id                       SERIAL PRIMARY KEY,
    first_name               VARCHAR(255),
    last_name                VARCHAR(255),
    chat_id                  BIGINT UNIQUE NOT NULL,
    username                 VARCHAR(255),
    created_at               TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMP     NOT NULL DEFAULT NOW(),
    published_at             TIMESTAMP,
    hidden_at                TIMESTAMP,
    notifications_enabled_at TIMESTAMP,
    avatar_url               VARCHAR(512),
    title                    VARCHAR(255),
    description              TEXT,
    language_code            VARCHAR(2)    NOT NULL DEFAULT 'en',
    country                  VARCHAR(255),
    city                     VARCHAR(255),
    country_code             VARCHAR(2),
    followers_count          INTEGER       NOT NULL DEFAULT 0,
    requests_count           INTEGER       NOT NULL DEFAULT 0
);

CREATE INDEX users_chat_id_index ON users (chat_id);

CREATE TABLE badges
(
    id         SERIAL PRIMARY KEY,
    text       VARCHAR(255) NOT NULL,
    icon       VARCHAR(255),
    color      VARCHAR(7),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE opportunities
(
    id          SERIAL PRIMARY KEY,
    text VARCHAR(255) NOT NULL,
    description TEXT,
    icon        VARCHAR(255),
    color       VARCHAR(7),
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE user_badges
(
    user_id  INTEGER REFERENCES users (id),
    badge_id INTEGER REFERENCES badges (id),
    UNIQUE (user_id, badge_id)
);

-- user_opportunities

CREATE TABLE user_opportunities
(
    user_id        INTEGER REFERENCES users (id),
    opportunity_id INTEGER REFERENCES opportunities (id),
    UNIQUE (user_id, opportunity_id)
);


CREATE TABLE collaborations
(
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER REFERENCES users (id),
    opportunity_id INTEGER REFERENCES opportunities (id),
    title          VARCHAR(255) NOT NULL,
    description    TEXT,
    is_payable     BOOLEAN      NOT NULL DEFAULT FALSE,
    published_at   TIMESTAMP,
    created_at     TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP    NOT NULL DEFAULT NOW(),
    country        VARCHAR(255),
    city           VARCHAR(255),
    country_code   VARCHAR(2),
    requests_count INTEGER      NOT NULL DEFAULT 0
);


CREATE TABLE user_followers
(
    user_id     INTEGER REFERENCES users (id),
    follower_id INTEGER REFERENCES users (id),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, follower_id)
);

CREATE OR REPLACE FUNCTION update_follower_counts() RETURNS TRIGGER AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE users SET followers_count = followers_count + 1 WHERE id = NEW.user_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE users SET followers_count = followers_count - 1 WHERE id = OLD.user_id;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_follower_counts_trigger
    AFTER INSERT OR DELETE
    ON user_followers
    FOR EACH ROW
EXECUTE FUNCTION update_follower_counts();


CREATE TYPE collaboration_request_status AS ENUM ('pending', 'approved', 'rejected');


CREATE TABLE collaboration_requests
(
    id               SERIAL PRIMARY KEY,
    collaboration_id INTEGER REFERENCES collaborations (id),
    user_id          INTEGER REFERENCES users (id),
    message          TEXT,
    created_at       TIMESTAMP                    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP                    NOT NULL DEFAULT NOW(),
    status           collaboration_request_status NOT NULL DEFAULT 'pending',
    UNIQUE (collaboration_id, user_id)
);

CREATE OR REPLACE FUNCTION update_collaboration_requests_count() RETURNS TRIGGER AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE collaborations SET requests_count = requests_count + 1 WHERE id = NEW.collaboration_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE collaborations SET requests_count = requests_count - 1 WHERE id = OLD.collaboration_id;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE user_collaboration_requests
(
    id           SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id),
    requester_id INTEGER REFERENCES users (id),
    message      TEXT,
    status       collaboration_request_status NOT NULL DEFAULT 'pending',
    created_at   TIMESTAMP                    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP                    NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, requester_id)
);

CREATE OR REPLACE FUNCTION update_user_requests_count() RETURNS TRIGGER AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE users SET requests_count = requests_count + 1 WHERE id = NEW.receiver_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE users SET requests_count = requests_count - 1 WHERE id = OLD.receiver_id;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;


CREATE table collaboration_badges
(
    collaboration_id INTEGER REFERENCES collaborations (id),
    badge_id         INTEGER REFERENCES badges (id),
    UNIQUE (collaboration_id, badge_id)
);

CREATE table notifications
(
    id         SERIAL PRIMARY KEY,
    user_id    BIGINT,
    message_id VARCHAR(255),
    chat_id    BIGINT,
    text       TEXT,
    image_url  TEXT,
    sent_at    TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

insert into badges (icon, color, text)
values ('f149', '#685155', 'Dog Father'),
       ('e561', '#17BEBB', 'Restaurateur'),
       ('f1e8', '#FE5F55', 'Wine Lover'),
       ('ea67', '#FF8C42', 'Career Specialist'),
       ('ea78', '#FE5F55', 'Yoga Lover'),
       ('e8d3', '#17BEBB', 'Mentor'),
       ('e521', '#FE5F55', 'Musician'),
       ('e905', '#FF8C42', 'Traveller'),
       ('f10a', '#17BEBB', 'Content Creator'),
       ('eb39', '#FF8C42', 'Founder'),
       ('e7c8', '#EF5DA8', 'Entrepreneur'),
       ('e06b', '#93961F', 'Product Designer'),
       ('f8d7', '#0B5351', 'Manager'),
       ('f84d', '#EF5DA8', 'Software Engineer'),
       ('e992', '#93961F', 'Business Developer'),
       ('e586', '#93961F', 'Architect'),
       ('ead1', '#3478F6', 'Data Scientist'),
       ('e273', '#EF5DA8', 'Fashion Designer'),
       ('f05b', '#3478F6', 'Marketer'),
       ('e4fb', '#0B5351', 'Investment Specialist'),
       ('ebc5', '#3478F6', 'FinTech'),
       ('eb43', '#17BEBB', 'Sportsman'),
       ('e87e', '#3478F6', 'Health Expert'),
       ('f221', '#FE5F55', 'BioTech Specialist'),
       ('ea53', '#685155', 'Foodie'),
       ('e80c', '#0B5351', 'Learner'),
       ('f882', '#685155', 'AI Engineer'),
       ('e9fe', '#FF8C42', '3D Modeller'),
       ('eb43', '#685155', 'Gym Rat'),
       ('ea34', '#FE5F55', 'Grappler'),
       ('eae9', '#17BEBB', 'Fighter'),
       ('e92c', '#3478F6', 'Lifelong Learner'),
       ('ebed', '#17BEBB', 'Medcine Doctor'),
       ('f345', '#3478F6', 'Writer'),
       ('f1ac', '#FF8C42', 'Teacher'),
       ('e875', '#685155', 'Back-end Developer'),
       ('ea19', '#FF8C42', 'Reader'),
       ('f1ea', '#3478F6', 'Fish Farmer'),
       ('ea79', '#0B5351', 'Farmer'),
       ('ea3a', '#0B5351', 'Biologist'),
       ('e760', '#93961F', 'Sustainability Expert'),
       ('e90e', '#EF5DA8', 'Lawyer'),
       ('ea30', '#685155', 'Sports Coach'),
       ('ea5f', '#17BEBB', 'Economist'),
       ('ea4b', '#93961F', 'Chemist'),
       ('e8f4', '#FE5F55', 'Visionary'),
       ('e9fe', '#FF8C42', '3D Artist'),
       ('f89a', '#685155', 'Failure'),
       ('e405', '#FE5F55', 'Music Producer'),
       ('e661', '#17BEBB', 'Motion Designer'),
       ('e2eb', '#FE5F55', 'Investor'),
       ('e3ae', '#FF8C42', 'Illustrator'),
       ('e32a', '#685155', 'Cyber Security'),
       ('e412', '#17BEBB', 'Photographer'),
       ('e3ae', '#FF8C42', 'Graphic Designer'),
       ('e86f', '#FE5F55', 'Front-end Developer'),
       ('f10a', '#17BEBB', 'UX Designer'),
       ('f10a', '#FF8C42', 'UI Designer'),
       ('e56c', '#FE5F55', 'Chef'),
       ('eff7', '#17BEBB', 'Taxi Driver'),
       ('e540', '#FF8C42', 'Bartender'),
       ('eb44', '#FE5F55', 'Barista'),
       ('e536', '#17BEBB', 'Walker'),
       ('e552', '#FF8C42', 'Pizza Lover'),
       ('ea53', '#FE5F55', 'Baker'),
       ('ea13', '#0B5351', 'CTO '),
       ('e3dd', '#3478F6', 'Cloud Engineer'),
       ('e9f9', '#0B5351', 'Golang'),
       ('e420', '#3478F6', 'Dancer'),
       ('f8d9', '#17BEBB', 'Project Manager'),
       ('eaf0', '#FE5F55', 'Business Assistant'),
       ('e52f', '#FF8C42', 'JavaScript'),
       ('e427', '#3478F6', 'Figma'),
       ('e320', '#EF5DA8', 'DevOps'),
       ('e91d', '#EF5DA8', 'Cat Mother'),
       ('e406', '#0B5351', 'Volunteer'),
       ('e3b7', '#3478F6', 'Web Designer'),
       ('e50a', '#FF8C42', 'Hiking'),
       ('ea4a', '#0B5351', 'Cloud Guru'),
       ('f1f3', '#FF8C42', 'Chiller'),
       ('f01f', '#17BEBB', 'Music Lover'),
       ('f06c', '#EF5DA8', 'Copilot enthusiast'),
       ('ea28', '#FF8C42', 'Gamer'),
       ('f05d', '#685155', 'Modeling'),
       ('ea1f', '#0B5351', 'Urbanist'),
       ('e52f', '#685155', 'bike rides'),
       ('e508', '#0B5351', 'Gardening'),
       ('e95f', '#FE5F55', 'Ghh'),
       ('f1b7', '#3478F6', 'Head of ESG'),
       ('e80b', '#3478F6', 'Sustainability Executive'),
       ('e566', '#3478F6', 'Runner'),
       ('e174', '#3478F6', 'Executive search'),
       ('e97a', '#EF5DA8', 'Recruitment Expert'),
       ('e853', '#FE5F55', 'Employee branding '),
       ('ea41', '#17BEBB', 'director of operations'),
       ('f041', '#17BEBB', 'Financial model'),
       ('eb1b', '#3478F6', 'Typescript'),
       ('e2bd', '#FF8C42', 'AWS'),
       ('e663', '#0B5351', 'VueJs'),
       ('e2be', '#3478F6', 'Google Cloud'),
       ('e0b2', '#0B5351', 'SEO'),
       ('e92c', '#EF5DA8', 'HRlover'),
       ('ea70', '#17BEBB', 'Psychotherapist'),
       ('e3ae', '#FE5F55', 'Painter'),
       ('e3b3', '#3478F6', 'Video'),
       ('e764', '#3478F6', 'Handmade'),
       ('ea75', '#3478F6', 'Sound design'),
       ('e953', '#FE5F55', 'Random'),
       ('e92c', '#0B5351', 'Zoomer');

insert into opportunities (icon, color, text, description)
values ('ea66', '#FF8C42', 'Acting in a play', 'Being an actor in a stage production'),
       ('e8da', '#3478F6', 'Acting in film', 'Being an actor in a movie or show'),
       ('ea1a', '#FF8C42', 'Activism', 'Supporting campaigns to bring about political change'),
       ('ea40', '#17BEBB', 'Advising companies', 'Advising companies on certain topics'),
       ('eb9b', '#FE5F55', 'Advising early stage companies', 'Providing insights to help a young company'),
       ('ebbc', '#0B5351', 'Advising late stage companies', 'Help mature companies get to the next level'),
       ('ef70', '#0B5351', 'Advising on SEO', 'Sharing SEO expertise'),
       ('e063', '#17BEBB', 'Appearing in music videos', 'Making a guest appearance in a music video'),
       ('e404', '#685155', 'Appearing in videos', 'Making a guest appearance in a video'),
       ('e43a', '#3478F6', 'Artist Management', 'Managing up and coming musicians'),
       ('ec0b', '#FE5F55', 'Beat Submissions', 'Listening to beats from up & coming producers'),
       ('e2eb', '#685155', 'Becoming a VC scout', 'Identifying and sharing new businesses that show potential'),
       ('f8d7', '#FF8C42', 'Being a brand ambassador', 'Promote a brand you''re passionate about'),
       ('e8da', '#17BEBB', 'Being an extra in a film or show', 'Appearing as an extra in a film or show'),
       ('e868', '#FE5F55', 'Beta testing new products', 'Checking out the newest consumer or business products'),
       ('e8af', '#3478F6', 'Brainstorming',
        'Joining video calls to brainstorm with like minded people on certain problems'),
       ('ebcb', '#FF8C42', 'Brand partnership', 'Partner with a brand for mutual benefit'),
       ('f184', '#FF8C42', 'Brand strategy consulting', 'Advise companies on their vision and brand'),
       ('e729', '#EF5DA8', 'Business partnerships', 'Connecting and building new business relationships'),
       ('f22e', '#17BEBB', 'Career coaching', 'Helping people make career decisions or prepare for interviews'),
       ('ea30', '#685155', 'Coaching founders', 'Providing direction and feedback to founders'),
       ('e7f0', '#3478F6', 'Co-founding a company', 'Building up a business from vision to reality with a co-founder'),
       ('e91f', '#685155', 'Consulting', 'Providing consultative services in your area of expertise'),
       ('e65f', '#EF5DA8', 'Content Creation', 'Creating content for brands and companies'),
       ('e746', '#0B5351', 'Copywriting', 'Writing for a product, marketing initiative, or project'),
       ('e7a3', '#FE5F55', 'Creating 2D animations', 'Rendering 2D visuals and animations'),
       ('e84d', '#0B5351', 'Creating 3D animations', 'Rendering 3D visuals and animations'),
       ('f1c1', '#FF8C42', 'Creating a design system', 'Designing a visual system for a brand or new concept'),
       ('e71c', '#FE5F55', 'Creating animations', 'Storytelling through a moving sequence of illustrations'),
       ('e3ae', '#EF5DA8', 'Creating illustrations', 'Designing a visual explanation for a concept or idea'),
       ('e25e', '#685155', 'Designing brands',
        'Creating a design identity for a new idea or reinventing an existing one'),
       ('f19e', '#3478F6', 'Designing clothing', 'Designing clothing for a brand or concept'),
       ('e545', '#0B5351', 'Designing floral arrangements', 'Arranging flowers for decoration'),
       ('e167', '#FE5F55', 'Designing fonts', 'Creating custom fonts or font systems'),
       ('e586', '#685155', 'Designing homes', 'Providing architectural designs for homes'),
       ('ead5', '#3478F6', 'Designing logos', 'Creating a visual that represents a brand''s identity'),
       ('eaf0', '#0B5351', 'Designing pitch decks',
        'Designing a visual aid that provides an overview of a company or startup'),
       ('eae4', '#685155', 'Designing products', 'Creating a product that provides a solution for users'' needs'),
       ('e06b', '#17BEBB', 'Designing Websites', 'Designing a website'),
       ('f10a', '#17BEBB', 'Design projects', 'Designing graphics or illustrations for a brand or new concept'),
       ('e1b0', '#0B5351', 'Developing apps', 'Writing code to bring an app design to life'),
       ('ef54', '#FE5F55', 'Developing a web application', 'Creating a new web application'),
       ('e338', '#EF5DA8', 'Developing games', 'Bringing a game concept to life'),
       ('e051', '#17BEBB', 'Developing sites in Webflow', 'Building sites witihin the Webflow platform'),
       ('ef54', '#FE5F55', 'Developing websites', 'Writing code and bringing a website design to life'),
       ('e666', '#0B5351', 'Editing books', 'Proofreading and critiquing book drafts'),
       ('e01f', '#FE5F55', 'Editing videos', 'Trim, edit and finalize video projects'),
       ('e0e6', '#3478F6', 'Email marketing consulting', 'Assisting teams with email strategy and execution'),
       ('ec09', '#FF8C42', 'Externships', 'Participating in a temporary job training program'),
       ('e404', '#FF8C42', 'Filming videos', 'Meet up and shoot video footage'),
       ('eaf8', '#EF5DA8', 'Fractional executive roles', 'Working part-time in an executive capacity'),
       ('e894', '#3478F6', 'Freelance roles', 'Freelancing, consulting, or part time work'),
       ('e943', '#EF5DA8', 'Full time roles', 'Starting a new role in a full time capacity'),
       ('ea70', '#3478F6', 'Fundraising for non-profits', 'Helping non-profits raise money for causes'),
       ('e865', '#FE5F55', 'Giving book recommendations', 'Sharing your favorite reads'),
       ('e86f', '#3478F6', 'Giving code reviews', 'Checking code for any errors + providing feedback'),
       ('e84f', '#685155', 'Giving college application feedback',
        'Providing feedback to help college applicants improve their applications'),
       ('f1c5', '#685155', 'Giving design crits', 'Providing constructive critiques on design work'),
       ('e560', '#FE5F55', 'Giving design feedback', 'Providing feedback on someone''s design'),
       ('e0c9', '#FF8C42', 'Giving feedback', 'Providing feedback to help someone improve'),
       ('e0c9', '#685155', 'Giving feedback on marketing copy',
        'Proofreading and giving constructive feedback on copy'),
       ('e2eb', '#EF5DA8', 'Giving fundraising advice',
        'Sharing insight and experience on successfully raising funding'),
       ('e1b8', '#0B5351', 'Giving music feedback', 'Listening to and giving feedback on pieces'),
       ('e521', '#EF5DA8', 'Giving music lessons',
        'Providing formal instruction to help people improve musical skills'),
       ('ea0e', '#FE5F55', 'Giving pitch deck feedback', 'Providing constructive critique on pitch decks'),
       ('e8eb', '#17BEBB', 'Giving portfolio feedback', 'Reviewing and giving feedback on someone''s portfolio'),
       ('e8eb', '#3478F6', 'Giving product feedback', 'Providing feedback to help improve a product'),
       ('e8cd', '#3478F6', 'Giving product reviews', 'Providing constructive feedback about products'),
       ('e85e', '#0B5351', 'Giving resume feedback',
        'Reviewing and critiquing the resumes of people in similar roles to me'),
       ('e80b', '#17BEBB', 'Giving travel advice', 'Sharing tips with other travelers'),
       ('e52f', '#0B5351', 'Going cycling', 'Going on bike rides with other cyclists'),
       ('e50a', '#FF8C42', 'Going hiking', 'Exploring mountains and trails with a group'),
       ('e566', '#3478F6', 'Going running', 'Meeting up for a casual run or training'),
       ('efef', '#0B5351', 'Grabbing a coffee', 'Meeting and talking over a cup of joe'),
       ('e663', '#0B5351', 'Graphic design', 'Creating digital visuals'),
       ('e31d', '#EF5DA8', 'Guest lecturing', 'Sharing your knowledge by guest lecturing at colleges and universities'),
       ('ea62', '#FF8C42', 'Having a dinner party', 'Getting together for a meal'),
       ('ea67', '#3478F6', 'Hiring', 'Hiring for full time roles'),
       ('e878', '#0B5351', 'Hosting events', 'Providing space for an event'),
       ('e58a', '#FF8C42', 'House swapping', 'Trading home bases and exploring new spaces'),
       ('e5d3', '#FF8C42', 'Instagram collaborations', 'Partnering to create Instagram content'),
       ('e80c', '#685155', 'Internships', 'Joining companies or organizations as an intern'),
       ('ef63', '#FE5F55', 'Investing', 'Investing in startups through crowdfunding and direct investments'),
       ('e019', '#FE5F55', 'Jam sessions', 'Making music with a group'),
       ('f8e0', '#FE5F55', 'Joining a band', 'Joining a group of musicians to create projects'),
       ('e98b', '#685155', 'Joining a book club', 'Reading and discussing books with a group'),
       ('f8d9', '#685155', 'Joining a community', 'Meeting other users with common interests'),
       ('efee', '#685155', 'Joining company boards', 'Sitting on the board of companies and nonprofits'),
       ('e0ef', '#EF5DA8', 'Joining DAOs', 'Joining Decentralized Autonomous Organizations'),
       ('e311', '#17BEBB', 'Joining Discords', 'Get in on a new Discord channel'),
       ('ea2f', '#17BEBB', 'Joining fantasy football leagues', 'Build your dream team and connect over sports'),
       ('ebbb', '#EF5DA8', 'Joining investment syndicates',
        'Joining a group of investors to pool capital together and fund deals'),
       ('e91d', '#EF5DA8', 'Joining parent groups', 'Meet, connect, and discuss with other parents'),
       ('eb9d', '#FF8C42', 'Just chatting', 'Chatting and networking with new people. Feel free to reach out!'),
       ('e55f', '#FF8C42', 'Live Streaming', 'Live streaming with people'),
       ('f041', '#685155', 'LP investing', 'Using limited partnerships for investing'),
       ('f0d3', '#FE5F55', 'Making angel investments', 'Invest first in ideas you''re passionate about'),
       ('f205', '#EF5DA8', 'Making early stage investments', 'Invest early in ideas you''re passionate about'),
       ('f8d7', '#17BEBB', 'Making Introductions', 'Introducing people to my network for potential opportunities'),
       ('f1d0', '#FF8C42', 'Making late stage investments', 'Take a growing company to the next level'),
       ('eb62', '#EF5DA8', 'Making real estate investments', 'Investing in real estate properties'),
       ('e04b', '#685155', 'Making videos', 'Creating and editing video content'),
       ('eb3d', '#FF8C42', 'Melody Loop Submissions', 'Listening to melody loops from up & coming producers'),
       ('e8d3', '#FF8C42', 'Mentoring', 'Mentoring people in similar fields to you'),
       ('e429', '#3478F6', 'Mixing and mastering', 'Combining musical tracks to create a new sound'),
       ('e83b', '#FE5F55', 'Mock interviewing', 'Practicing interview skills'),
       ('e8fc', '#3478F6', 'Modeling', 'Modeling for photoshoots'),
       ('ef3d', '#3478F6', 'Moderating clubs', 'Taking on a leadership role in a club'),
       ('e8af', '#FF8C42', 'Moderating Discord channels', 'Upholding community guidelines for a Discord channel'),
       ('e939', '#3478F6', 'Moderating events', 'Running an event as a moderator or emcee'),
       ('e32e', '#3478F6', 'Musical collaborations', 'Collaborating with other musicians to produce a musical project'),
       ('e7f0', '#FF8C42', 'Networking', 'Meeting with new people to expand your connections'),
       ('e992', '#17BEBB', 'New business leads', 'Sourcing new clients'),
       ('e0d0', '#FE5F55', 'New Client Inquiries', 'Reviewing new client requests and inquiries'),
       ('e943', '#17BEBB', 'New roles', 'Providing new job opportunities'),
       ('e7f9', '#3478F6', 'NFT projects', 'Partnering with people on NFT projects'),
       ('f1b8', '#3478F6', 'Open Source Contributions', 'Contributing to open source projects'),
       ('eb8e', '#FF8C42', 'Organizing a hackathon', 'Putting together a hackathon event'),
       ('ebcc', '#FF8C42', 'Organizing events', 'Planning events'),
       ('e4fb', '#FE5F55', 'Organizing marketing campaigns', 'Planning and executing marketing campaigns'),
       ('ead1', '#0B5351', 'Participating in hackathons',
        'Collaborating with other developers to hack together programs'),
       ('efd7', '#EF5DA8', 'Participating in User Research', 'Sharing your thoughts on early products or services'),
       ('ebcb', '#3478F6', 'Partnering on Side Projects', 'Partnering with people to build side projects part time'),
       ('e422', '#0B5351', 'Part time roles', 'Freelancing, consulting, or part time employment'),
       ('efec', '#FE5F55', 'Peer editing', 'Providing input and feedback on educational projects'),
       ('ea66', '#FF8C42', 'Performing', 'Entertaining an audience through performance'),
       ('e92c', '#FF8C42', 'Performing dance', 'Dancing for a live audience'),
       ('e610', '#3478F6', 'Performing music', 'Presenting pieces of music for a live audience'),
       ('e8c8', '#FF8C42', 'Performing standup comedy', 'Performing comedy for a live audience'),
       ('e850', '#FF8C42', 'Philanthropy', 'Donating large sums of money to organizations in need of support'),
       ('e8fc', '#3478F6', 'Photography', 'Taking photoshoots'),
       ('e770', '#685155', 'Planning digital ad campaigns', 'Mapping out digital advertising and promotion'),
       ('ea26', '#17BEBB', 'Playing basketball', 'Meeting up for pickup games or shooting hoops'),
       ('e9b0', '#FF8C42', 'Playing chess', 'Exercise your strategy muscles with a game of chess'),
       ('ea28', '#EF5DA8', 'Playing games', 'Have fun and connect over games'),
       ('eb31', '#FF8C42', 'Playing Pokemon Go', 'Gotta catch em all!'),
       ('ea2f', '#EF5DA8', 'Playing soccer', 'Meet up and connect over the game of soccer'),
       ('ea33', '#17BEBB', 'Playing sports', 'Have fun and connect over sports'),
       ('eb40', '#17BEBB', 'Playing tabletop games', 'Have fun and connect over board games'),
       ('ea32', '#17BEBB', 'Playing tennis', 'Meet up and connect over tennis'),
       ('ea28', '#EF5DA8', 'Playing video games', 'Have fun and connect over video games'),
       ('e03d', '#685155', 'Playlisting', 'Creating playlists'),
       ('e8e2', '#3478F6', 'Practicing a new language', 'Improving your language skills through conversation'),
       ('e927', '#EF5DA8', 'Practicing Arabic', 'Improving your Arabic skills through conversation'),
       ('e927', '#0B5351', 'Practicing Chinese', 'Improving your Chinese skills through conversation'),
       ('e927', '#17BEBB', 'Practicing English', 'Improving your English skills through conversation'),
       ('e927', '#EF5DA8', 'Practicing French', 'Improving your French skills through conversation'),
       ('e927', '#FF8C42', 'Practicing German', 'Improving your German skills through conversation'),
       ('e927', '#EF5DA8', 'Practicing Hindi', 'Improving your Hindi skills through conversation'),
       ('e927', '#685155', 'Practicing Indonesian', 'Improving your Indonesian skills through conversation'),
       ('e927', '#0B5351', 'Practicing Italian', 'Improving your Italian skills through conversation'),
       ('e927', '#FE5F55', 'Practicing Korean', 'Improving your Korean skills through conversation'),
       ('e927', '#685155', 'Practicing Norwegian', 'Improving your Norwegian skills through conversation'),
       ('e927', '#EF5DA8', 'Practicing Portuguese', 'Improving your Portuguese skills through conversation'),
       ('e927', '#FE5F55', 'Practicing Russian', 'Improving your Russian skills through conversation'),
       ('e927', '#17BEBB', 'Practicing Spanish', 'Improving your Spanish skills through conversation'),
       ('ea78', '#FE5F55', 'Practicing yoga', 'Joining a yoga session. Namaste!'),
       ('f048', '#17BEBB', 'Producing a Podcast', 'Producing a podcast in a freelance capacity'),
       ('e8da', '#FF8C42', 'Producing a short film', 'Creating and editing content for a short film'),
       ('e880', '#685155', 'Proofreading', 'Help review a book or article before it''s published'),
       ('efed', '#685155', 'Providing interior design services', 'Designing and curating a space'),
       ('ef56', '#EF5DA8', 'Providing legal aid', 'Assisting people who are unable to afford legal representation'),
       ('e8a1', '#685155', 'Providing sponsorships', 'Giving financial support to people or organizations'),
       ('e8e2', '#FF8C42', 'Providing translations', 'Translating content between languages'),
       ('f1b6', '#FF8C42', 'Raising funding', 'Fundraising for your company'),
       ('e4fc', '#FE5F55', 'Research projects', 'Creating a project to answer a research question'),
       ('e666', '#FE5F55', 'Reviewing book proposals', 'Giving constructive feedback on book ideas'),
       ('f11b', '#3478F6', 'Running statistical analyses', 'Providing data insights or analysis'),
       ('ea3e', '#EF5DA8', 'Sharing lesson plans', 'Share your content with other educators'),
       ('eb81', '#0B5351', 'Social media management', 'Manage social channels and engagement'),
       ('eb53', '#EF5DA8', 'Social takeovers', 'Running someone''s social account for a short time'),
       ('e86f', '#FF8C42', 'Software development projects', 'Connect and assist on coding projects'),
       ('e27c', '#0B5351', 'Speaking at Events', 'Speaking on topics you are passionate about at all kinds of events'),
       ('e83b', '#FF8C42', 'Speaking on Clubhouse', 'Speaking on certain topics or hosting a Clubhouse show'),
       ('e83b', '#0B5351', 'Speaking on podcasts', 'Speaking on podcast episodes on certain topics'),
       ('e0e5', '#FF8C42', 'Speaking on Twitter Spaces', 'Speaking on certain topics or hosting a Twitter Space'),
       ('e83b', '#0B5351', 'Sponsoring content', 'Paying for promotion through sponsored content'),
       ('ea65', '#0B5351', 'Sponsoring events', 'Paying for promotion through an event'),
       ('f18a', '#3478F6', 'Sponsoring newsletters', 'Paying for promotion through a newsletter'),
       ('ea14', '#17BEBB', 'Starting new open source projects',
        'Creating a new project for open source contributors to join'),
       ('e333', '#FF8C42', 'Streaming on Twitch', 'Livestreaming on the Twitch platform'),
       ('f06a', '#17BEBB', 'Streaming on YouTube', 'Livestream on the YouTube platform'),
       ('e80c', '#FF8C42', 'Student projects', 'Forming or joining a student project'),
       ('e32d', '#17BEBB', 'Studio Sessions', 'Open to Sessions'),
       ('e666', '#0B5351', 'Study sessions', 'Forming or joining a study group'),
       ('ea9b', '#685155', 'Talent scouting', 'Searching for artists for employment'),
       ('eb81', '#17BEBB', 'Talking to journalists', 'Giving quotes to journalists on specific topics'),
       ('f8ea', '#3478F6', 'Teaching', 'Sharing area expertise or interest with learners'),
       ('f06c', '#0B5351', 'Teaching AI', 'Giving lessons on AI'),
       ('f10a', '#3478F6', 'Teaching design', 'Sharing design expertise with learners'),
       ('e26b', '#EF5DA8', 'Teaching entrepreneurship', 'Sharing entrepreneurship insights with learners'),
       ('eb8e', '#3478F6', 'Teaching software engineering', 'Sharing software engineering expertise with learners'),
       ('e41b', '#17BEBB', 'TikTok collaborations', 'Partnering to create TikTok content'),
       ('efec', '#685155', 'Tutoring', 'Teaching and sharing knowledge on area of expertise'),
       ('f1ab', '#3478F6', 'Volunteering', 'Paying it forward by giving up your time to help organizations in need'),
       ('e7e9', '#685155', 'Wedding planning', 'Help organize the wedding plan'),
       ('f88c', '#0B5351', 'Writing', 'Contributing articles or blogs on certain topics'),
       ('e745', '#3478F6', 'Writing blog posts', 'Create new content for a topic you''re interest in'),
       ('f85a', '#17BEBB', 'Youtube collaborations', 'Partnering to create YouTube content');

