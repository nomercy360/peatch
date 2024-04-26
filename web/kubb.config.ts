import { defineConfig } from '@kubb/core';
import { definePlugin as createSwagger } from '@kubb/swagger';
import { definePlugin as createSwaggerTS } from '@kubb/swagger-ts';

export default defineConfig(async () => {
  return {
    root: '.',
    input: {
      path: '../docs/swagger.yaml',
    },
    output: {
      path: './gen',
    },
    plugins: [
      createSwagger(
        {
          output: false,
          validate: true,
        },
      ),
      createSwaggerTS(
        {},
      ),
    ],
  };
});