FROM scratch
LABEL Name=ssm_run
LABEL Author=davyj0nes

ADD ssm_run_static /ssm_run

ENTRYPOINT ["./ssm_run"]
