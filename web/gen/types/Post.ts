import type { UserProfile } from "./UserProfile";

 export type Post = {
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
    hidden_at?: string;
    /**
     * @type integer | undefined
    */
    id?: number;
    /**
     * @type string | undefined
    */
    image_url?: string;
    /**
     * @type string | undefined
    */
    title?: string;
    /**
     * @type string | undefined
    */
    updated_at?: string;
    user?: UserProfile;
    /**
     * @type integer | undefined
    */
    user_id?: number;
};