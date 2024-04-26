import type { User } from "./User";

 export type GetApiUsersChatIdPathParams = {
    /**
     * @description Chat ID
     * @type integer
    */
    chat_id: number;
};

 /**
 * @description OK
*/
export type GetApiUsersChatId200 = User;

 /**
 * @description OK
*/
export type GetApiUsersChatIdQueryResponse = User;

 export type GetApiUsersChatIdQuery = {
    Response: GetApiUsersChatIdQueryResponse;
    PathParams: GetApiUsersChatIdPathParams;
};