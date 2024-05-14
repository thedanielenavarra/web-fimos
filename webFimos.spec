Name:       webFimos
Version:    1.0
Release:    1%{?dist}
Summary:    webFimos installation
License:    MIT

# Define dependencies if needed
Requires:   systemd golang

# Define the source files
Source0:    webFimos
Source1:    web-fimos.service

# Define the installation paths
%define _systemdunitdir /usr/lib/systemd/system

%description
Your description here

%prep

%build

%install
mkdir -p %{buildroot}/etc/web-fimos
mkdir -p %{buildroot}/var/log/web-fimos
mkdir -p %{buildroot}/usr/bin
install -m 755 %{SOURCE0} %{buildroot}/usr/bin/webFimos

mkdir -p %{buildroot}%{_systemdunitdir}
install -m 644 %{SOURCE1} %{buildroot}%{_systemdunitdir}/web-fimos.service

%post
/usr/bin/systemctl daemon-reload
webFimos /etc/web-fimos/web-fimos.json --regen --host localhost --port 9091 > /var/log/web-fimos/web-fimos.log 2>&1

%files
%{_systemdunitdir}/web-fimos.service
/usr/bin/webFimos
/etc/web-fimos/
/var/log/web-fimos/
%changelog

