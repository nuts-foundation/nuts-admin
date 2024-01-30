export default {
  install: (app, apiOptions = {}) => {
    const { defaultOptions } = apiOptions

    const authHeader = () => {
      const sessionToken = localStorage.getItem('session')
      if (sessionToken) {
        return { Authorization: `Bearer ${sessionToken}` }
      }
      return {}
    }
    const api = {}

    const httpMethods = ['get', 'post', 'put', 'delete']
    httpMethods.forEach((method) => {
      api[method] = (url, data = null, requestOptions = {}) => {
        const options = {
          ...defaultOptions,
          method: method.toUpperCase(),
          headers: {
            'Content-Type': 'application/json',
            ...authHeader()
          },
          ...requestOptions
        }
        if (data) {
          options.body = JSON.stringify(data)
        }

        return fetch(url, options)
          .then((response) => {
            // Parse JSON or just grab plain text
            const contentType = response.headers.get("content-type");
            var parsedResponse = undefined
            var isJson = false
            if (contentType && contentType.indexOf("application/json") !== -1) {
              parsedResponse = response.json()
              isJson = true
            } else {
              parsedResponse = response.text()
            }

            return parsedResponse
              .then((data) => {
                if (response.ok) {
                  return Promise.resolve(data)
                } else {
                  if (apiOptions.forbiddenRoute && response.status === 401) {
                    return app.config.globalProperties.$router.push(apiOptions.forbiddenRoute)
                  } else {
                    if (isJson) {
                      return Promise.reject(data.error)
                    } else {
                      return Promise.reject(data)
                    }
                  }
                }
              }).catch(reason => {
                // Handle 404 since it does not have content and the response.json() will fail.
                if (response.status === 404) {
                  return Promise.reject(response.statusText)
                }
                // Handle 204 since it does not have content and the response.json() will fail.
                if (response.status === 204) {
                  return Promise.resolve(response)
                }
                // Handle 201 since it might not have content and the response.json() will fail.
                if (response.status === 201) {
                  return Promise.resolve(response)
                }
                return Promise.reject(reason)
              })
          })
      }
    })

    app.config.globalProperties.$api = api
  }
}
