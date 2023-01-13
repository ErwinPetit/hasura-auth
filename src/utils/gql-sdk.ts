// TODO get rid of graphql-request
import { GraphQLClient } from 'graphql-request';
import { ENV } from './env';

export const client = new GraphQLClient(ENV.HASURA_GRAPHQL_GRAPHQL_URL, {
  headers: {
    'x-hasura-admin-secret': ENV.HASURA_GRAPHQL_ADMIN_SECRET,
  },
});
