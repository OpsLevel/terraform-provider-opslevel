mock_resource "opslevel_integration_google_cloud" {
  defaults = {
    ownership_tag_keys = []
    projects = [
      {
        id   = "pauls-test-project-431117"
        name = "Paul's Test Project"
        url  = "https://console.cloud.google.com/home/dashboard?project=pauls-test-project-431117"
      },
      {
        id   = "powerful-surf-427415"
        name = "Powerful Surf"
        url  = "https://console.cloud.google.com/home/dashboard?project=powerful-surf-427415"
      }
    ]
  }
}

