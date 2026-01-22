/**
 * Direct access to the current Go process environment and metadata.
 * Mimics a subset of the Node.js process API.
 */
declare const process: {
    /**
     * Environment variables. Only variables whitelisted or prefixed with TYPEGO_ 
     * are accessible for security.
     */
    env: Record<string, string>;

    /**
     * Operating system platform (e.g., 'windows', 'linux', 'darwin').
     */
    platform: string;

    /**
     * Returns the current working directory.
     */
    cwd(): string;

    /**
     * Array of command-line arguments.
     */
    argv: string[];

    /**
     * Go runtime version.
     */
    version: string;
};
