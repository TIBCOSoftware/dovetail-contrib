###################################################################################
#                          contribs/dovetail                                      #
###################################################################################
FROM scratch

VOLUME [ /var/lib/dovetail/dovetail-contrib ]

COPY . /var/lib/dovetail/dovetail-contrib/