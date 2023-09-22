# resource "installer_apt" "neofetch" {
#     name = "neofetch"
# }

# resource "installer_apt" "example" {
#   configurable_attribute = "some-value"

#   connection {
#       host = "hi"
#       user = "root"
#       private_key = "mom"
#       agent = false
#       timeout = "2m"
#   }
# }

resource "installer_script" "test" {
  path           = "/tmp/installer-myapp-test"
  install_script = <<-EOF
  /bin/bash

  touch /tmp/installer-myapp-test
  chmod +x /tmp/installer-myapp-test
  exit 0
  EOF

  uninstall_script = <<-EOF
  /bin/bash

  rm -f /tmp/installer-myapp-test
  exit 0
  EOF
}
