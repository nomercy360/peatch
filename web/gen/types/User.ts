import type { Badge } from "./Badge";
import type { Opportunity } from "./Opportunity";

 export type User = {
    /**
     * @type string | undefined
    */
    avatar_url?: string;
    /**
     * @type array | undefined
    */
    badges?: Badge[];
    /**
     * @type integer | undefined
    */
    chat_id?: number;
    /**
     * @type string | undefined
    */
    city?: string;
    /**
     * @type string | undefined
    */
    country?: string;
    /**
     * @type string | undefined
    */
    country_code?: string;
    /**
     * @type string | undefined
    */
    created_at?: string;
    /**
     * @type string | undefined
    */
    description?: string;
    /**
     * @type string | undefined
    */
    first_name?: string;
    /**
     * @type integer | undefined
    */
    followers_count?: number;
    /**
     * @type integer | undefined
    */
    following_count?: number;
    /**
     * @type string | undefined
    */
    hidden_at?: string;
    /**
     * @type integer | undefined
    */
    id?: number;
    /**
     * @type boolean | undefined
    */
    is_following?: boolean;
    /**
     * @type boolean | undefined
    */
    is_liked?: boolean;
    /**
     * @type string | undefined
    */
    last_check_in?: string;
    /**
     * @type string | undefined
    */
    last_name?: string;
    /**
     * @type integer | undefined
    */
    likes_count?: number;
    /**
     * @type array | undefined
    */
    opportunities?: Opportunity[];
    /**
     * @type integer | undefined
    */
    peatch_points?: number;
    /**
     * @type string | undefined
    */
    published_at?: string;
    /**
     * @type string | undefined
    */
    title?: string;
    /**
     * @type string | undefined
    */
    username?: string;
};