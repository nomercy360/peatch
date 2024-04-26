import type { UserWithToken } from './UserWithToken';

export type GetApiAuthTelegramQueryParams = {
    /**
     * @description Query ID
     * @type string
    */
    query_id: string;
    /**
     * @description User
     * @type string
    */
    user: string;
    /**
     * @description Auth date
     * @type string
    */
    auth_date: string;
    /**
     * @description Hash
     * @type string
    */
    hash: string;
};

 /**
 * @description OK
*/
 export type GetApiAuthTelegram200 = UserWithToken;

 /**
 * @description OK
*/
 export type GetApiAuthTelegramQueryResponse = UserWithToken;

 export type GetApiAuthTelegramQuery = {
    Response: GetApiAuthTelegramQueryResponse;
    QueryParams: GetApiAuthTelegramQueryParams;
};