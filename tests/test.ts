import 'mocha';
import { request, spec } from 'pactum';
import { faker } from '@faker-js/faker';
import { Badge, CreateCollaboration, UpdateUserRequest } from '../web/gen';

const baseUrl = 'http://localhost:8080/api';
const authUrl = 'http://localhost:8080/auth/telegram';

describe('Test Admin Routes', () => {
  before(async () => {
    request.setDefaultTimeout(10000);
  });

  const tgAuthQuery = 'query_id=AAH9mUo3AAAAAP2ZSjec547F&user=%7B%22id%22%3A927635965%2C%22first_name%22%3A%22Maksim%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22mkkksim%22%2C%22language_code%22%3A%22en%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&auth_date=1714101898&hash=b9dca58f25c765f33bf2651fa25b44b8ff38e33b400951cdd3ed68ea4c024af7';
  const tgSecondUser = 'query_id=AAH9mUo3AAAAAP2ZSjec547F&user=%7B%22id%22%3A123456789%2C%22first_name%22%3A%22John%22%2C%22last_name%22%3A%22Doe%22%2C%22username%22%3A%22johndoe%22%2C%22language_code%22%3A%22en%22%2C%22is_premium%22%3Afalse%2C%22allows_write_to_pm%22%3Afalse%7D&auth_date=1714101898&hash=767e3daef07fc1f3ee89b421f316f9b2f0a1b8997a7a73c2df68e2dea17ce644';
  const tgThirdUser = 'query_id=AAH9mUo3AAAAAP2ZSjec547F&user=%7B%22id%22%3A987654321%2C%22first_name%22%3A%22Fabio%22%2C%22last_name%22%3A%22Rossi%22%2C%22username%22%3A%22fabiorossi%22%2C%22language_code%22%3A%22it%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&auth_date=1714101898&hash=e602a16800c36af4fcc77ccd6b1032d54c3cecd9fe4ea1b52cb10b5bd431e4eb';

  const users = { 'firstUser': tgAuthQuery, 'secondUser': tgSecondUser, 'thirdUser': tgThirdUser };

  for (const [name, query] of Object.entries(users)) {
    it('should authenticate ' + name, async () => {
      await spec()
        .post(`${authUrl}?${query}`)
        .expectStatus(200)
        .expectJsonSchema({
          type: 'object',
          required: ['token', 'user'],
        })
        .stores(name + 'Token', 'token');
    });
  }

  const newBadge: Badge = {
    text: faker.word.noun({ length: { min: 5, max: 20 } }) + ' badge',
    icon: 'f1f9',
    color: 'ff0000',
  };

  const titleCase = (str: string) => {
    // first letter of each word to uppercase
    return str.replace(/\w\S*/g, (txt) => {
      return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
    });
  };

  const firstUserAuth = { 'Authorization': 'Bearer $S{firstUserToken}' };
  const secondUserAuth = { 'Authorization': 'Bearer $S{secondUserToken}' };

  it('should create badge', async () => {
    await spec()
      .post(`${baseUrl}/badges`)
      .withHeaders(firstUserAuth)
      .withJson(newBadge)
      .expectStatus(201)
      .expectJsonMatch({
        text: titleCase(newBadge.text!),
        icon: newBadge.icon,
        color: newBadge.color,
      })
      .expectJsonSchema({
        type: 'object',
        required: ['id', 'text', 'icon', 'created_at', 'color'],
      })
      .stores('badgeId', 'id');
  });

  it('should retrieve badge', async () => {
    await spec()
      .get(`${baseUrl}/badges`)
      .withHeaders(firstUserAuth)
      .expectStatus(200)
      .stores('firstBadgeId', '[0].id')
      .stores('secondBadgeId', '[1].id')
      .stores('thirdBadgeId', '[2].id')
      .stores('fourthBadgeId', '[3].id')
      .stores('fifthBadgeId', '[4].id')
      .stores('sixthBadgeId', '[5].id')
      .stores('seventhBadgeId', '[6].id')
      .stores('eighthBadgeId', '[7].id')
      .stores('ninthBadgeId', '[8].id')
      .stores('tenthBadgeId', '[9].id')
      .stores('eleventhBadgeId', '[10].id')
      .stores('twelfthBadgeId', '[11].id');
  });

  it('should retrieve opportunities', async () => {
    await spec()
      .get(`${baseUrl}/opportunities`)
      .withHeaders(firstUserAuth)
      .expectStatus(200)
      .expectJsonSchema({
        type: 'array',
        items: {
          type: 'object',
          required: ['id', 'text', 'description', 'icon', 'color', 'created_at'],
        },
      })
      .stores('firstOpportunityId', '[0].id')
      .stores('secondOpportunityId', '[1].id')
      .stores('thirdOpportunityId', '[2].id');
  });

  const firstUserProfile: UpdateUserRequest = {
    first_name: faker.person.firstName(),
    last_name: faker.person.lastName(),
    title: faker.person.jobTitle(),
    description: faker.person.jobDescriptor(),
    city: faker.location.city(),
    country: faker.location.country(),
    country_code: faker.location.countryCode(),
    avatar_url: faker.image.avatar(),
    //@ts-ignore
    opportunity_ids: ['$S{firstOpportunityId}'],
    //@ts-ignore
    badge_ids: ['$S{firstBadgeId}', '$S{secondBadgeId}', '$S{thirdBadgeId}', '$S{fourthBadgeId}', '$S{fifthBadgeId}'],
  };

  const secondUserProfile: UpdateUserRequest = {
    first_name: faker.person.firstName(),
    last_name: faker.person.lastName(),
    title: faker.person.jobTitle(),
    description: faker.person.jobDescriptor(),
    city: faker.location.city(),
    country: faker.location.country(),
    country_code: faker.location.countryCode(),
    avatar_url: faker.image.avatar(),
    //@ts-ignore
    opportunity_ids: ['$S{firstOpportunityId}', '$S{thirdOpportunityId}'],
    //@ts-ignore
    badge_ids: ['$S{firstBadgeId}', '$S{secondBadgeId}'],
  };

  const userUpdateProfiles = {
    'firstUser': { auth: 'firstUserToken', profile: firstUserProfile },
    'secondUser': { auth: 'secondUserToken', profile: secondUserProfile },
  };

  for (const [name, { auth, profile }] of Object.entries(userUpdateProfiles)) {
    it('should updates user profile ' + name, async () => {
      await spec()
        .put(`${baseUrl}/users`)
        .withHeaders({
          Authorization: 'Bearer $S{' + auth + '}',
        })
        .withJson(profile)
        .expectStatus(200)
        .expectJsonMatch({
          first_name: profile.first_name,
          last_name: profile.last_name,
          title: profile.title,
          description: profile.description,
          city: profile.city,
          country: profile.country,
          country_code: profile.country_code,
          avatar_url: profile.avatar_url,
        })
        .expectJsonSchema({
          type: 'object',
          required: ['id', 'first_name', 'last_name', 'title', 'description', 'city', 'country', 'country_code', 'avatar_url', 'created_at', 'updated_at'],
        })
        .stores(name + 'Id', 'id');
    });

    it('should get' + name + ' user profile', async () => {
      await spec()
        .get(`${baseUrl}/users/$S{${name}Id}`)
        .withHeaders({
          'Authorization': 'Bearer $S{' + auth + '}',
        })
        .expectStatus(200)
        .expectJsonSchema({
          type: 'object',
          required: ['id', 'first_name', 'last_name', 'title', 'description', 'city', 'country', 'country_code', 'avatar_url', 'created_at', 'updated_at'],
        })
        .expectJsonMatch({
          id: '$S{' + name + 'Id}',
          first_name: profile.first_name,
          last_name: profile.last_name,
          title: profile.title,
          description: profile.description,
          city: profile.city,
          country: profile.country,
          country_code: profile.country_code,
          avatar_url: profile.avatar_url,
        })
        .expectJsonLength('badges', profile.badge_ids.length)
        .expectJsonLength('opportunities', profile.opportunity_ids.length);
    });
  }

  it('should list users', async () => {
    await spec()
      .get(`${baseUrl}/users`)
      .withHeaders(firstUserAuth)
      .expectStatus(200)
      .expectJsonSchema({
        type: 'array',
        items: {
          type: 'object',
          required: ['id', 'first_name', 'last_name', 'title', 'description', 'city', 'country', 'country_code', 'avatar_url', 'created_at', 'updated_at'],
        },
      })
      .expectJsonLength(0);
  });

  it('should publish user', async () => {
    await spec()
      .post(`${baseUrl}/users/publish`)
      .withHeaders(firstUserAuth)
      .expectStatus(204);
  });

  it('should list published user', async () => {
    await spec()
      .get(`${baseUrl}/users`)
      .withHeaders(firstUserAuth)
      .expectStatus(200)
      .expectJsonLength(1);
  });

  it('should hide user', async () => {
    await spec()
      .post(`${baseUrl}/users/hide`)
      .withHeaders(firstUserAuth)
      .expectStatus(204);
  });

  it('should not show hidden user', async () => {
    await spec()
      .get(`${baseUrl}/users/$S{firstUserId}`)
      .withHeaders(secondUserAuth)
      .expectStatus(404);
  });

  it('should follow user', async () => {
    await spec()
      .post(`${baseUrl}/users/$S{firstUserId}/follow`)
      .withHeaders(secondUserAuth)
      .expectStatus(204);
  });


  it('should get user follower count', async () => {
    await spec()
      .get(`${baseUrl}/users/$S{firstUserId}`)
      .withHeaders(firstUserAuth)
      .expectStatus(200)
      .expectJsonMatch({
        followers_count: 1,
      });
  });


  it('should unfollow user', async () => {
    await spec()
      .delete(`${baseUrl}/users/$S{firstUserId}/follow`)
      .withHeaders(secondUserAuth)
      .expectStatus(204);
  });

  it('should get user follower count', async () => {
    await spec()
      .get(`${baseUrl}/users/$S{firstUserId}`)
      .withHeaders(firstUserAuth)
      .expectStatus(200)
      .expectJsonMatch({
        followers_count: 0,
      });
  });

  const firstCollaboration: CreateCollaboration = {
    badge_ids: [],
    city: faker.location.city(),
    country: faker.location.country(),
    country_code: faker.location.countryCode(),
    description: faker.lorem.sentence(),
    is_payable: false,
    opportunity_id: 0,
    title: faker.lorem.sentence(),
  };

  const secondCollaboration: CreateCollaboration = {
    badge_ids: [],
    city: faker.location.city(),
    country: faker.location.country(),
    country_code: faker.location.countryCode(),
    description: faker.lorem.sentence(),
    is_payable: false,
    opportunity_id: 0,
    title: faker.lorem.sentence(),
  };

  it('user can create collaboration', async () => {
    await spec()
      .post(`${baseUrl}/collaborations`)
      .withHeaders(firstUserAuth)
      .withJson({
        ...firstCollaboration,
        opportunity_id: '$S{firstOpportunityId}',
        badge_ids: ['$S{firstBadgeId}', '$S{secondBadgeId}', '$S{thirdBadgeId}', '$S{fourthBadgeId}', '$S{fifthBadgeId}', '$S{sixthBadgeId}', '$S{seventhBadgeId}', '$S{eighthBadgeId}', '$S{ninthBadgeId}', '$S{tenthBadgeId}'],
      })
      .expectStatus(201)
      .expectJsonSchema({
        type: 'object',
        required: ['id', 'title', 'description', 'city', 'country', 'country_code', 'is_payable', 'created_at', 'updated_at'],
      })
      .stores('firstCollaborationId', 'id');
  });

  it('user can create second collaboration', async () => {
    await spec()
      .post(`${baseUrl}/collaborations`)
      .withHeaders(firstUserAuth)
      .withJson({
        ...secondCollaboration,
        opportunity_id: '$S{secondOpportunityId}',
        badge_ids: ['$S{eleventhBadgeId}', '$S{twelfthBadgeId}'],
      })
      .expectStatus(201)
      .expectJsonSchema({
        type: 'object',
        required: ['id', 'title', 'description', 'city', 'country', 'country_code', 'is_payable', 'created_at', 'updated_at'],
      })
      .stores('secondCollaborationId', 'id');
  });

  it('should list collaborations', async () => {
    await spec()
      .get(`${baseUrl}/collaborations`)
      .withHeaders(secondUserAuth)
      .expectStatus(200)
      .expectJsonSchema({
        type: 'array',
        items: {
          type: 'object',
          required: ['id', 'title', 'description', 'city', 'country', 'country_code', 'is_payable', 'created_at', 'updated_at', 'opportunity'],
        },
      })
      .expectJsonLength(0);
  });

  it('should publish collaboration', async () => {
    await spec()
      .post(`${baseUrl}/collaborations/$S{firstCollaborationId}/publish`)
      .withHeaders(firstUserAuth);
  });

  it('should publish second collaboration', async () => {
    await spec()
      .post(`${baseUrl}/collaborations/$S{secondCollaborationId}/publish`)
      .withHeaders(firstUserAuth);
  });

  it('should list published collaboration', async () => {
    await spec()
      .get(`${baseUrl}/collaborations`)
      .withHeaders(secondUserAuth)
      .expectStatus(200)
      .expectJsonLength(2);
  });

  it('should get collaboration', async () => {
    await spec()
      .get(`${baseUrl}/collaborations/$S{firstCollaborationId}`)
      .withHeaders(firstUserAuth)
      .expectStatus(200)
      .expectJsonSchema({
        type: 'object',
        required: ['id', 'title', 'description', 'city', 'country', 'country_code', 'is_payable', 'created_at', 'updated_at', 'opportunity', 'badges'],
      })
      .expectJsonMatch({
        id: '$S{firstCollaborationId}',
        title: firstCollaboration.title,
        description: firstCollaboration.description,
        city: firstCollaboration.city,
        country: firstCollaboration.country,
        country_code: firstCollaboration.country_code,
        is_payable: firstCollaboration.is_payable,
        opportunity: {
          id: '$S{firstOpportunityId}',
        },
      })
      .expectJsonLength('badges', 10);
  });

  it('should get collaboration', async () => {
    await spec()
      .get(`${baseUrl}/collaborations/$S{secondCollaborationId}`)
      .withHeaders({ 'Authorization': 'Bearer $S{firstUserToken}' })
      .expectStatus(200)
      .expectJsonSchema({
        type: 'object',
        required: ['id', 'title', 'description', 'city', 'country', 'country_code', 'is_payable', 'created_at', 'updated_at', 'opportunity', 'badges'],
      })
      .expectJsonMatch({
        id: '$S{secondCollaborationId}',
        title: secondCollaboration.title,
        description: secondCollaboration.description,
        city: secondCollaboration.city,
        country: secondCollaboration.country,
        country_code: secondCollaboration.country_code,
        is_payable: secondCollaboration.is_payable,
        opportunity: {
          id: '$S{secondOpportunityId}',
        },
        user: {
          id: '$S{firstUserId}',
        },
      })
      .expectJsonLength('badges', 2);
  });

  it('should hide collaboration', async () => {
    await spec()
      .post(`${baseUrl}/collaborations/$S{secondCollaborationId}/hide`)
      .withHeaders(firstUserAuth)
      .expectStatus(204);
  });

  it('should not show hidden collaboration', async () => {
    await spec()
      .get(`${baseUrl}/collaborations/$S{secondCollaborationId}`)
      .withHeaders(secondUserAuth)
      .expectStatus(404);
  });

  it('should search collaborations', async () => {
    await spec()
      .get(`${baseUrl}/collaborations?search=${firstCollaboration.title}`)
      .withHeaders(firstUserAuth)
      .expectStatus(200)
      .expectJsonLength(1);
  });

  it('should search collaborations', async () => {
    await spec()
      .get(`${baseUrl}/collaborations?search=${secondCollaboration.title}`)
      .withHeaders(firstUserAuth)
      .expectStatus(200)
      .expectJsonLength(1);
  });
});