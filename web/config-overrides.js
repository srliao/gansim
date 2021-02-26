const {
    override, 
    getBabelLoader, 
    addWebpackModuleRule
  } = require('customize-cra');
  
  module.exports = (config, env) => {
    const babelLoader = getBabelLoader(config);
    
    return override(
      addWebpackModuleRule({
        test: /\.mjs$/,
        include: /node_modules/,
        type: "javascript/auto",
      })
    )(config, env)
  }