resource "installer_apt" "neofetch" {
    name = "neofetch"
    #sudo = false

    # remote_connection {
    #     type = "ssh"
    #     host = "changeme"
    #     user = "root"
    #     password = "changeme"
    # }
}

# resource "installer_apt" "dos2unix" {
#     depends_on = [installer_apt.neofetch]
#     name = "dos2unix"
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

# resource "installer_script" "test" {
#   path           = "/tmp/installer-myapp-test"
#   install_script = <<-EOF
#   echo INSTALLING!
#   touch /tmp/installer-myapp-test
#   chmod +x /tmp/installer-myapp-test
#   exit 0
#   EOF

#   uninstall_script = <<-EOF
#   echo UNINSTALLING!
#   rm -f /tmp/installer-myapp-test
#   exit 0
#   EOF

#   #shell = "bash"
#   depends_on = [installer_apt.neofetch]
# }
