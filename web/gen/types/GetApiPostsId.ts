import type { Post } from "./Post";

 export type GetApiPostsIdPathParams = {
    /**
     * @description Post ID
     * @type integer
    */
    id: number;
};

 /**
 * @description OK
*/
export type GetApiPostsId200 = Post;

 /**
 * @description OK
*/
export type GetApiPostsIdQueryResponse = Post;

 export type GetApiPostsIdQuery = {
    Response: GetApiPostsIdQueryResponse;
    PathParams: GetApiPostsIdPathParams;
};